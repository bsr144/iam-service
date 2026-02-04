package internal

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CompleteRegistration sets password, creates user account, and optionally returns auth tokens.
// Design reference: .claude/doc/email-otp-signup-api.md Section 2.6
func (uc *usecase) CompleteRegistration(
	ctx context.Context,
	tenantID, registrationID uuid.UUID,
	registrationToken string,
	req *authdto.CompleteRegistrationRequest,
	ipAddress, userAgent string,
) (*authdto.CompleteRegistrationResponse, error) {
	// 1. Validate registration token
	_, err := uc.validateRegistrationCompleteToken(registrationToken, registrationID, tenantID)
	if err != nil {
		return nil, err
	}

	// 2. Get registration session (returns AppError directly)
	session, err := uc.Redis.GetRegistrationSession(ctx, tenantID, registrationID)
	if err != nil {
		return nil, err
	}

	// 3. Check session status
	if session.IsExpired() {
		return nil, errors.New("REGISTRATION_EXPIRED", "Registration session has expired", http.StatusGone)
	}

	if session.Status != entity.RegistrationSessionStatusVerified {
		return nil, errors.ErrForbidden("Email has not been verified")
	}

	// 4. Verify token hash matches (single-use check)
	tokenHash := sha256.Sum256([]byte(registrationToken))
	tokenHashStr := hex.EncodeToString(tokenHash[:])
	if session.RegistrationTokenHash == nil || *session.RegistrationTokenHash != tokenHashStr {
		return nil, errors.ErrUnauthorized("Registration token has already been used or is invalid")
	}

	// 5. Validate password
	if err := uc.validatePassword(req.Password); err != nil {
		return nil, errors.ErrValidation(err.Error())
	}

	// 6. Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternal("failed to hash password").WithError(err)
	}
	passwordHashStr := string(passwordHash)

	// 7. Check if email was registered in the meantime (race condition protection)
	emailExists, err := uc.UserRepo.EmailExistsInTenant(ctx, tenantID, session.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to check email").WithError(err)
	}
	if emailExists {
		return nil, errors.ErrConflict("This email has already been registered")
	}

	// 8. Determine user status (for now, always active - tenant approval logic can be added later)
	userStatus := entity.UserStatusActive
	requiresApproval := false

	// TODO: Check tenant settings for requires_approval flag when TenantSettings is loaded
	// For now, default to no approval required

	// 9. Create user in transaction
	var userID uuid.UUID
	now := time.Now()

	err = uc.DB.Transaction(func(tx *gorm.DB) error {
		userID = uuid.New()

		// Create user
		user := &entity.User{
			UserID:        userID,
			TenantID:      &tenantID,
			Email:         session.Email,
			EmailVerified: true, // Already verified via OTP
			IsActive:      userStatus == entity.UserStatusActive,
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// Create credentials
		credentials := &entity.UserCredentials{
			UserCredentialID: uuid.New(),
			UserID:           userID,
			PasswordHash:     &passwordHashStr,
			PasswordHistory:  json.RawMessage("[]"),
			PINHistory:       json.RawMessage("[]"),
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		if err := tx.Create(credentials).Error; err != nil {
			return err
		}

		// Create profile
		profile := &entity.UserProfile{
			UserProfileID: uuid.New(),
			UserID:        userID,
			FirstName:     req.FirstName,
			LastName:      req.LastName,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if req.PhoneNumber != nil {
			profile.Phone = req.PhoneNumber
		}
		if err := tx.Create(profile).Error; err != nil {
			return err
		}

		// Create security state
		security := &entity.UserSecurity{
			UserSecurityID: uuid.New(),
			UserID:         userID,
			Metadata:       json.RawMessage("{}"),
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := tx.Create(security).Error; err != nil {
			return err
		}

		// Create activation tracking
		tracking := entity.NewUserActivationTracking(userID, &tenantID)
		if err := tracking.AddStatusTransition(string(userStatus), "registration"); err != nil {
			return err
		}
		if err := tx.Create(tracking).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, errors.ErrInternal("failed to create user").WithError(err)
	}

	// 10. Cleanup Redis session
	_ = uc.Redis.DeleteRegistrationSession(ctx, tenantID, registrationID)
	_ = uc.Redis.UnlockRegistrationEmail(ctx, tenantID, session.Email)

	// 11. Build response
	response := &authdto.CompleteRegistrationResponse{
		UserID: userID,
		Email:  session.Email,
		Status: string(userStatus),
		Profile: authdto.RegistrationUserProfile{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		},
	}

	// 12. Auto-login if user is active (not pending approval)
	if !requiresApproval {
		response.Message = "Registration completed successfully. You are now logged in."

		// Generate access and refresh tokens
		accessToken, refreshToken, expiresIn, err := uc.generateAuthTokensForRegistration(ctx, userID, tenantID, session.Email)
		if err != nil {
			// Don't fail registration if token generation fails
			// User can still login manually
			response.Message = "Registration completed successfully. Please login to continue."
		} else {
			tokenType := "Bearer"
			response.AccessToken = &accessToken
			response.RefreshToken = &refreshToken
			response.TokenType = &tokenType
			response.ExpiresIn = &expiresIn
		}

		// Send welcome email
		_ = uc.EmailService.SendWelcome(ctx, session.Email, req.FirstName)
	} else {
		response.Message = "Registration submitted. Your account is pending administrator approval. You will receive an email when your account is activated."
	}

	return response, nil
}

// generateAuthTokensForRegistration generates access and refresh tokens for auto-login after registration.
// TODO: This should integrate with your existing token generation infrastructure.
func (uc *usecase) generateAuthTokensForRegistration(ctx context.Context, userID, tenantID uuid.UUID, email string) (string, string, int, error) {
	// For now, return an error to indicate this needs implementation
	// In production, this would:
	// 1. Generate access token with proper claims (user_id, tenant_id, roles, auth_method="registration")
	// 2. Generate refresh token
	// 3. Store refresh token in the database
	// 4. Return tokens with expiry

	// Placeholder - this causes the response to say "Please login to continue"
	return "", "", 0, errors.ErrInternal("token generation pending implementation")
}

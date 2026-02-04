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

func (uc *usecase) CompleteRegistration(
	ctx context.Context,
	tenantID, registrationID uuid.UUID,
	registrationToken string,
	req *authdto.CompleteRegistrationRequest,
	ipAddress, userAgent string,
) (*authdto.CompleteRegistrationResponse, error) {
	_, err := uc.validateRegistrationCompleteToken(registrationToken, registrationID, tenantID)
	if err != nil {
		return nil, err
	}

	session, err := uc.Redis.GetRegistrationSession(ctx, tenantID, registrationID)
	if err != nil {
		return nil, err
	}

	if session.IsExpired() {
		return nil, errors.New("REGISTRATION_EXPIRED", "Registration session has expired", http.StatusGone)
	}

	if session.Status != entity.RegistrationSessionStatusVerified {
		return nil, errors.ErrForbidden("Email has not been verified")
	}

	tokenHash := sha256.Sum256([]byte(registrationToken))
	tokenHashStr := hex.EncodeToString(tokenHash[:])
	if session.RegistrationTokenHash == nil || *session.RegistrationTokenHash != tokenHashStr {
		return nil, errors.ErrUnauthorized("Registration token has already been used or is invalid")
	}

	if err := uc.validatePassword(req.Password); err != nil {
		return nil, errors.ErrValidation(err.Error())
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternal("failed to hash password").WithError(err)
	}
	passwordHashStr := string(passwordHash)

	emailExists, err := uc.UserRepo.EmailExistsInTenant(ctx, tenantID, session.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to check email").WithError(err)
	}
	if emailExists {
		return nil, errors.ErrConflict("This email has already been registered")
	}

	userStatus := entity.UserStatusActive
	requiresApproval := false

	var userID uuid.UUID
	now := time.Now()

	err = uc.DB.Transaction(func(tx *gorm.DB) error {
		userID = uuid.New()

		user := &entity.User{
			UserID:        userID,
			TenantID:      &tenantID,
			Email:         session.Email,
			EmailVerified: true,
			IsActive:      userStatus == entity.UserStatusActive,
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}

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

	_ = uc.Redis.DeleteRegistrationSession(ctx, tenantID, registrationID)
	_ = uc.Redis.UnlockRegistrationEmail(ctx, tenantID, session.Email)

	response := &authdto.CompleteRegistrationResponse{
		UserID: userID,
		Email:  session.Email,
		Status: string(userStatus),
		Profile: authdto.RegistrationUserProfile{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		},
	}

	if !requiresApproval {
		response.Message = "Registration completed successfully. You are now logged in."

		accessToken, refreshToken, expiresIn, err := uc.generateAuthTokensForRegistration(ctx, userID, tenantID, session.Email)
		if err != nil {
			response.Message = "Registration completed successfully. Please login to continue."
		} else {
			tokenType := "Bearer"
			response.AccessToken = &accessToken
			response.RefreshToken = &refreshToken
			response.TokenType = &tokenType
			response.ExpiresIn = &expiresIn
		}

		_ = uc.EmailService.SendWelcome(ctx, session.Email, req.FirstName)
	} else {
		response.Message = "Registration submitted. Your account is pending administrator approval. You will receive an email when your account is activated."
	}

	return response, nil
}

func (uc *usecase) generateAuthTokensForRegistration(ctx context.Context, userID, tenantID uuid.UUID, email string) (string, string, int, error) {

	return "", "", 0, errors.ErrInternal("token generation pending implementation")
}

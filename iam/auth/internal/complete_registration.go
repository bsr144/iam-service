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

	var user *entity.User
	now := time.Now()

	err = uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		user = &entity.User{
			TenantID:      &tenantID,
			Email:         session.Email,
			EmailVerified: true,
			IsActive:      userStatus == entity.UserStatusActive,
		}
		if err := uc.UserRepo.Create(txCtx, user); err != nil {
			return err
		}

		credentials := &entity.UserCredentials{
			UserID:          user.ID,
			PasswordHash:    &passwordHashStr,
			PasswordHistory: json.RawMessage("[]"),
			PINHistory:      json.RawMessage("[]"),
			CreatedAt:       now,
			UpdatedAt:       now,
		}
		if err := uc.UserCredentialsRepo.Create(txCtx, credentials); err != nil {
			return err
		}

		profile := &entity.UserProfile{
			UserID:    user.ID,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if req.PhoneNumber != nil {
			profile.Phone = req.PhoneNumber
		}
		if err := uc.UserProfileRepo.Create(txCtx, profile); err != nil {
			return err
		}

		security := &entity.UserSecurity{
			UserID:    user.ID,
			Metadata:  json.RawMessage("{}"),
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := uc.UserSecurityRepo.Create(txCtx, security); err != nil {
			return err
		}

		tracking := entity.NewUserActivationTracking(user.ID, &tenantID)
		if err := tracking.AddStatusTransition(string(userStatus), "registration"); err != nil {
			return err
		}
		if err := uc.UserActivationTrackingRepo.Create(txCtx, tracking); err != nil {
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
		UserID: user.ID,
		Email:  session.Email,
		Status: string(userStatus),
		Profile: authdto.RegistrationUserProfile{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		},
	}

	if !requiresApproval {
		response.Message = "Registration completed successfully. You are now logged in."

		accessToken, refreshToken, expiresIn, err := uc.generateAuthTokensForRegistration(ctx, user.ID, tenantID, session.Email)
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

func (uc *usecase) generateAuthTokensForRegistration(_ context.Context, _, _ uuid.UUID, _ string) (string, string, int, error) {

	return "", "", 0, errors.ErrInternal("token generation pending implementation")
}

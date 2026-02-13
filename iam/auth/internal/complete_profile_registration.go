package internal

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"
)

func (uc *usecase) CompleteProfileRegistration(
	ctx context.Context,
	req *authdto.CompleteProfileRegistrationRequest,
) (*authdto.CompleteProfileRegistrationResponse, error) {
	_, err := uc.validateRegistrationCompleteToken(req.RegistrationToken, req.RegistrationID)
	if err != nil {
		return nil, err
	}

	session, err := uc.Redis.GetRegistrationSession(ctx, req.RegistrationID)
	if err != nil {
		return nil, err
	}

	if session.IsExpired() {
		return nil, errors.New("REGISTRATION_EXPIRED", "Registration session has expired", http.StatusGone)
	}

	if !session.IsPasswordSet() {
		return nil, errors.ErrForbidden("Registration session is not ready for profile completion. Password must be set first.")
	}

	tokenHash := sha256.Sum256([]byte(req.RegistrationToken))
	tokenHashStr := hex.EncodeToString(tokenHash[:])
	if session.RegistrationTokenHash == nil || *session.RegistrationTokenHash != tokenHashStr {
		return nil, errors.ErrUnauthorized("Registration token has already been used or is invalid")
	}

	if err := uc.validateProfileFields(req); err != nil {
		return nil, err
	}

	firstName, lastName := splitFullName(req.FullName)

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, errors.ErrValidation("Invalid date_of_birth format. Use YYYY-MM-DD")
	}

	age := calculateAge(dob)
	if age < 18 {
		return nil, errors.ErrValidation("You must be at least 18 years old to register")
	}

	dobStr := req.DateOfBirth
	gender := entity.Gender(req.Gender)
	maritalStatus := entity.MaritalStatus(req.MaritalStatus)

	emailExists, err := uc.UserRepo.EmailExists(ctx, session.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to check email").WithError(err)
	}
	if emailExists {
		return nil, errors.ErrConflict("This email has already been registered")
	}

	passwordHashStr, err := uc.Redis.GetRegistrationPasswordHash(ctx, req.RegistrationID)
	if err != nil {
		return nil, errors.ErrForbidden("Password has not been set")
	}

	userStatus := entity.UserStatusActive
	var user *entity.User
	now := time.Now()

	err = uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		user = &entity.User{
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
			UserID:        user.ID,
			FirstName:     firstName,
			LastName:      lastName,
			Phone:         &req.PhoneNumber,
			DateOfBirth:   &dobStr,
			Gender:        &gender,
			MaritalStatus: &maritalStatus,
			Address:       &req.Address,
			PlaceOfBirth:  &req.PlaceOfBirth,
			CreatedAt:     now,
			UpdatedAt:     now,
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

		tracking := entity.NewUserActivationTracking(user.ID, nil)
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

	_ = uc.Redis.DeleteRegistrationSession(ctx, req.RegistrationID)
	_ = uc.Redis.UnlockRegistrationEmail(ctx, session.Email)

	accessToken, refreshToken, expiresIn, err := uc.generateAuthTokensForRegistration(ctx, user.ID, session.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate auth tokens").WithError(err)
	}

	_ = uc.EmailService.SendWelcome(ctx, session.Email, firstName)

	response := &authdto.CompleteProfileRegistrationResponse{
		UserID:  user.ID,
		Email:   session.Email,
		Status:  string(userStatus),
		Message: "Registration completed successfully. You are now logged in.",
		Profile: authdto.RegistrationUserProfile{
			FirstName: firstName,
			LastName:  lastName,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}

	return response, nil
}

func (uc *usecase) validateProfileFields(req *authdto.CompleteProfileRegistrationRequest) error {
	if strings.TrimSpace(req.FullName) == "" {
		return errors.ErrValidation("full_name is required")
	}

	if req.PhoneNumber == "" {
		return errors.ErrValidation("phone_number is required")
	}

	if req.DateOfBirth == "" {
		return errors.ErrValidation("date_of_birth is required")
	}

	validGenders := map[string]bool{"male": true, "female": true, "other": true}
	if !validGenders[req.Gender] {
		return errors.ErrValidation("gender must be one of: male, female, other")
	}

	validMaritalStatuses := map[string]bool{"single": true, "married": true, "divorced": true, "widowed": true}
	if !validMaritalStatuses[req.MaritalStatus] {
		return errors.ErrValidation("marital_status must be one of: single, married, divorced, widowed")
	}

	if len(req.Address) < 10 {
		return errors.ErrValidation("address must be at least 10 characters")
	}

	if len(req.PlaceOfBirth) < 2 {
		return errors.ErrValidation("place_of_birth must be at least 2 characters")
	}

	return nil
}

func splitFullName(fullName string) (firstName, lastName string) {
	trimmed := strings.TrimSpace(fullName)
	lastSpaceIndex := strings.LastIndex(trimmed, " ")

	if lastSpaceIndex == -1 {
		return trimmed, ""
	}

	return strings.TrimSpace(trimmed[:lastSpaceIndex]), strings.TrimSpace(trimmed[lastSpaceIndex+1:])
}

func calculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()

	if now.YearDay() < birthDate.YearDay() {
		age--
	}

	return age
}

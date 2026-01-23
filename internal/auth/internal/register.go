package internal

import (
	"context"
	"encoding/json"
	"iam-service/entity"
	"iam-service/internal/auth/authdto"
	"iam-service/pkg/errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (uc *usecase) Register(ctx context.Context, req *authdto.RegisterRequest) (*authdto.RegisterResponse, error) {
	tenantExists, err := uc.TenantRepo.Exists(ctx, req.TenantID)
	if err != nil {
		return nil, errors.ErrInternal("failed to verify tenant").WithError(err)
	}
	if !tenantExists {
		return nil, errors.ErrTenantNotFound()
	}

	emailExists, err := uc.UserRepo.EmailExistsInTenant(ctx, req.TenantID, req.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to check email").WithError(err)
	}
	if emailExists {
		return nil, errors.ErrUserAlreadyExists()
	}

	if err := uc.validatePassword(req.Password); err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternal("failed to hash password").WithError(err)
	}

	otp, otpHash, err := uc.generateOTP()
	if err != nil {
		return nil, errors.ErrInternal("failed to generate OTP").WithError(err)
	}

	var response *authdto.RegisterResponse
	now := time.Now()
	passwordHashStr := string(passwordHash)

	err = uc.DB.Transaction(func(tx *gorm.DB) error {
		userID := uuid.New()
		user := &entity.User{
			UserID:        userID,
			TenantID:      &req.TenantID,
			Email:         req.Email,
			EmailVerified: false,
			IsActive:      true,
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

		tracking := entity.NewUserActivationTracking(userID, &req.TenantID)
		if err := tracking.AddStatusTransition(string(entity.UserStatusPendingOTPVerification), "system"); err != nil {
			return err
		}
		if err := tx.Create(tracking).Error; err != nil {
			return err
		}

		otpExpiry := now.Add(time.Duration(OTPExpiryMinutes) * time.Minute)
		verification := &entity.EmailVerification{
			EmailVerificationID: uuid.New(),
			TenantID:            req.TenantID,
			UserID:              userID,
			Email:               req.Email,
			OTPCode:             otp,
			OTPHash:             otpHash,
			OTPType:             entity.OTPTypeRegistration,
			ExpiresAt:           otpExpiry,
			CreatedAt:           now,
		}
		if err := tx.Create(verification).Error; err != nil {
			return err
		}

		response = &authdto.RegisterResponse{
			UserID:       userID,
			Email:        req.Email,
			Status:       string(entity.UserStatusPendingOTPVerification),
			OTPExpiresAt: otpExpiry,
		}

		return nil
	})

	if err != nil {
		return nil, errors.ErrInternal("failed to create user").WithError(err)
	}

	if err := uc.EmailService.SendOTP(ctx, req.Email, otp, OTPExpiryMinutes); err != nil {

	}

	return response, nil
}

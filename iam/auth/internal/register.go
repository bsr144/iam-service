package internal

import (
	"context"
	"encoding/json"
	"time"

	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"golang.org/x/crypto/bcrypt"
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

	err = uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		user := &entity.User{
			TenantID:      &req.TenantID,
			Email:         req.Email,
			EmailVerified: false,
			IsActive:      true,
		}
		if err := uc.UserRepo.Create(txCtx, user); err != nil {
			return err
		}

		credentials := &entity.UserCredentials{
			UserID:          user.UserID,
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
			UserID:    user.UserID,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := uc.UserProfileRepo.Create(txCtx, profile); err != nil {
			return err
		}

		security := &entity.UserSecurity{
			UserID:    user.UserID,
			Metadata:  json.RawMessage("{}"),
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := uc.UserSecurityRepo.Create(txCtx, security); err != nil {
			return err
		}

		tracking := entity.NewUserActivationTracking(user.UserID, &req.TenantID)
		if err := tracking.AddStatusTransition(string(entity.UserStatusPendingOTPVerification), "system"); err != nil {
			return err
		}
		if err := uc.UserActivationTrackingRepo.Create(txCtx, tracking); err != nil {
			return err
		}

		otpExpiry := now.Add(time.Duration(OTPExpiryMinutes) * time.Minute)
		verification := &entity.EmailVerification{
			TenantID:  req.TenantID,
			UserID:    user.UserID,
			Email:     req.Email,
			OTPCode:   otp,
			OTPHash:   otpHash,
			OTPType:   entity.OTPTypeRegistration,
			ExpiresAt: otpExpiry,
			CreatedAt: now,
		}
		if err := uc.EmailVerificationRepo.Create(txCtx, verification); err != nil {
			return err
		}

		response = &authdto.RegisterResponse{
			UserID:       user.UserID,
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

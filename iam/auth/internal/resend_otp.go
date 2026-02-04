package internal

import (
	"context"
	stderrors "errors"
	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"
	"time"

	"github.com/google/uuid"
)

func (uc *usecase) ResendOTP(ctx context.Context, req *authdto.ResendOTPRequest) (*authdto.ResendOTPResponse, error) {
	user, err := uc.UserRepo.GetByEmail(ctx, req.TenantID, req.Email)
	if err != nil {
		if stderrors.Is(err, errors.SentinelNotFound) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, errors.ErrInternal("failed to get user").WithError(err)
	}

	if user.EmailVerified {
		return nil, errors.ErrValidation("Email is already verified")
	}

	activeCount, err := uc.EmailVerificationRepo.CountActiveOTPsByEmail(ctx, req.Email, entity.OTPTypeRegistration)
	if err != nil {
		return nil, errors.ErrInternal("failed to check OTP count").WithError(err)
	}
	if activeCount >= MaxActiveOTPPerEmail {
		return nil, errors.ErrTooManyRequests("Too many active OTPs. Please wait before requesting a new one.")
	}

	otp, otpHash, err := uc.generateOTP()
	if err != nil {
		return nil, errors.ErrInternal("failed to generate OTP").WithError(err)
	}

	now := time.Now()
	otpExpiry := now.Add(time.Duration(OTPExpiryMinutes) * time.Minute)

	verification := &entity.EmailVerification{
		EmailVerificationID: uuid.New(),
		TenantID:            req.TenantID,
		UserID:              user.UserID,
		Email:               req.Email,
		OTPCode:             otp,
		OTPHash:             otpHash,
		OTPType:             entity.OTPTypeRegistration,
		ExpiresAt:           otpExpiry,
		CreatedAt:           now,
	}

	if err := uc.EmailVerificationRepo.Create(ctx, verification); err != nil {
		return nil, errors.ErrInternal("failed to create verification").WithError(err)
	}

	if err := uc.EmailService.SendOTP(ctx, req.Email, otp, OTPExpiryMinutes); err != nil {

	}

	return &authdto.ResendOTPResponse{
		OTPExpiresAt: otpExpiry,
	}, nil
}

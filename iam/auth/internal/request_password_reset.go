package internal

import (
	"context"
	stderrors "errors"
	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (uc *usecase) RequestPasswordReset(ctx context.Context, req *authdto.RequestPasswordResetRequest) (*authdto.RequestPasswordResetResponse, error) {
	user, err := uc.UserRepo.GetByEmail(ctx, req.TenantID, req.Email)
	if err != nil {
		if stderrors.Is(err, errors.SentinelNotFound) {
			// Return success to prevent user enumeration
			return &authdto.RequestPasswordResetResponse{
				OTPExpiresAt: time.Now().Add(time.Duration(OTPExpiryMinutes) * time.Minute),
				EmailMasked:  maskEmail(req.Email),
			}, nil
		}
		return nil, errors.ErrInternal("failed to get user").WithError(err)
	}

	if !user.IsActive {

		return &authdto.RequestPasswordResetResponse{
			OTPExpiresAt: time.Now().Add(time.Duration(OTPExpiryMinutes) * time.Minute),
			EmailMasked:  maskEmail(req.Email),
		}, nil
	}

	activeCount, err := uc.EmailVerificationRepo.CountActiveOTPsByEmail(ctx, req.Email, entity.OTPTypePasswordReset)
	if err != nil {
		return nil, errors.ErrInternal("failed to check rate limit").WithError(err)
	}
	if activeCount >= MaxActiveOTPPerEmail {
		return nil, errors.ErrTooManyRequests("Too many password reset requests. Please try again later.")
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
		OTPType:             entity.OTPTypePasswordReset,
		ExpiresAt:           otpExpiry,
		CreatedAt:           now,
	}

	if err := uc.EmailVerificationRepo.Create(ctx, verification); err != nil {
		return nil, errors.ErrInternal("failed to create verification record").WithError(err)
	}

	if err := uc.EmailService.SendOTP(ctx, req.Email, otp, OTPExpiryMinutes); err != nil {

	}

	return &authdto.RequestPasswordResetResponse{
		OTPExpiresAt: otpExpiry,
		EmailMasked:  maskEmail(req.Email),
	}, nil
}
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	local := parts[0]
	domain := parts[1]

	if len(local) <= 1 {
		return email
	}

	masked := string(local[0]) + strings.Repeat("*", len(local)-1) + "@" + domain
	return masked
}

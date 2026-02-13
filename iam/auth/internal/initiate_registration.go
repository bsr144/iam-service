package internal

import (
	"context"
	"time"

	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"
	"iam-service/pkg/logger"

	"github.com/google/uuid"
)

func (uc *usecase) InitiateRegistration(
	ctx context.Context,
	req *authdto.InitiateRegistrationRequest,
) (*authdto.InitiateRegistrationResponse, error) {
	emailExists, err := uc.UserRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to check email").WithError(err)
	}
	if emailExists {
		return &authdto.InitiateRegistrationResponse{
			RegistrationID: uuid.New().String(),
			Email:          req.Email,
			Status:         string(entity.RegistrationSessionStatusPendingVerification),
			Message:        "Verification code sent to your email",
			ExpiresAt:      time.Now().Add(time.Duration(RegistrationSessionExpiryMinutes) * time.Minute),
			OTPConfig: authdto.OTPConfig{
				Length:                RegistrationOTPLength,
				ExpiresInMinutes:      RegistrationOTPExpiryMinutes,
				MaxAttempts:           RegistrationOTPMaxAttempts,
				ResendCooldownSeconds: RegistrationOTPResendCooldown,
				MaxResends:            RegistrationOTPMaxResends,
			},
		}, nil
	}

	rateLimitTTL := time.Duration(RegistrationRateLimitWindow) * time.Minute
	count, err := uc.Redis.IncrementRegistrationRateLimit(ctx, req.Email, rateLimitTTL)
	if err != nil {
		return nil, err
	}
	if count > int64(RegistrationRateLimitPerHour) {
		return nil, errors.ErrTooManyRequests("Too many registration attempts. Please try again later.")
	}

	emailLocked, err := uc.Redis.IsRegistrationEmailLocked(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if emailLocked {
		return nil, errors.ErrConflict("An active registration already exists for this email")
	}

	otp, otpHash, err := uc.generateOTP()
	if err != nil {
		return nil, errors.ErrInternal("failed to generate OTP").WithError(err)
	}

	now := time.Now()
	sessionID := uuid.New()
	sessionTTL := time.Duration(RegistrationSessionExpiryMinutes) * time.Minute
	otpExpiry := now.Add(time.Duration(RegistrationOTPExpiryMinutes) * time.Minute)

	session := &entity.RegistrationSession{
		ID:                    sessionID,
		Email:                 req.Email,
		Status:                entity.RegistrationSessionStatusPendingVerification,
		OTPHash:               otpHash,
		OTPCreatedAt:          now,
		OTPExpiresAt:          otpExpiry,
		Attempts:              0,
		MaxAttempts:           RegistrationOTPMaxAttempts,
		ResendCount:           0,
		MaxResends:            RegistrationOTPMaxResends,
		ResendCooldownSeconds: RegistrationOTPResendCooldown,
		IPAddress:             req.IPAddress,
		UserAgent:             req.UserAgent,
		CreatedAt:             now,
		ExpiresAt:             now.Add(sessionTTL),
	}

	locked, err := uc.Redis.LockRegistrationEmail(ctx, req.Email, sessionTTL)
	if err != nil {
		return nil, errors.ErrInternal("failed to lock email").WithError(err)
	}
	if !locked {
		return nil, errors.ErrConflict("An active registration already exists for this email")
	}

	if err := uc.Redis.CreateRegistrationSession(ctx, session, sessionTTL); err != nil {

		_ = uc.Redis.UnlockRegistrationEmail(ctx, req.Email)
		return nil, err
	}

	if err := uc.EmailService.SendOTP(ctx, req.Email, otp, RegistrationOTPExpiryMinutes); err != nil {
		uc.AuditLogger.Log(ctx, logger.AuditEvent{
			Domain:  "auth",
			Action:  "registration_otp_send_failed",
			Success: false,
			Reason:  err.Error(),
			Metadata: map[string]any{
				"email": req.Email,
			},
		})
	}

	return &authdto.InitiateRegistrationResponse{
		RegistrationID: sessionID.String(),
		Email:          req.Email,
		Status:         string(entity.RegistrationSessionStatusPendingVerification),
		Message:        "Verification code sent to your email",
		ExpiresAt:      session.ExpiresAt,
		OTPConfig: authdto.OTPConfig{
			Length:                RegistrationOTPLength,
			ExpiresInMinutes:      RegistrationOTPExpiryMinutes,
			MaxAttempts:           RegistrationOTPMaxAttempts,
			ResendCooldownSeconds: RegistrationOTPResendCooldown,
			MaxResends:            RegistrationOTPMaxResends,
		},
	}, nil
}

package internal

import (
	"context"
	"time"

	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) InitiateRegistration(
	ctx context.Context,
	tenantID uuid.UUID,
	req *authdto.InitiateRegistrationRequest,
	ipAddress, userAgent string,
) (*authdto.InitiateRegistrationResponse, error) {
	tenant, err := uc.TenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrTenantNotFound()
		}
		return nil, errors.ErrInternal("failed to verify tenant").WithError(err)
	}
	if !tenant.IsActive() {
		return nil, errors.ErrTenantInactive()
	}

	emailExists, err := uc.UserRepo.EmailExistsInTenant(ctx, tenantID, req.Email)
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
	count, err := uc.Redis.IncrementRegistrationRateLimit(ctx, tenantID, req.Email, rateLimitTTL)
	if err != nil {
		return nil, err
	}
	if count > int64(RegistrationRateLimitPerHour) {
		return nil, errors.ErrTooManyRequests("Too many registration attempts. Please try again later.")
	}

	emailLocked, err := uc.Redis.IsRegistrationEmailLocked(ctx, tenantID, req.Email)
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
		TenantID:              tenantID,
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
		IPAddress:             ipAddress,
		UserAgent:             userAgent,
		CreatedAt:             now,
		ExpiresAt:             now.Add(sessionTTL),
	}

	locked, err := uc.Redis.LockRegistrationEmail(ctx, tenantID, req.Email, sessionTTL)
	if err != nil {
		return nil, errors.ErrInternal("failed to lock email").WithError(err)
	}
	if !locked {
		return nil, errors.ErrConflict("An active registration already exists for this email")
	}

	if err := uc.Redis.CreateRegistrationSession(ctx, session, sessionTTL); err != nil {

		_ = uc.Redis.UnlockRegistrationEmail(ctx, tenantID, req.Email)
		return nil, err
	}

	_ = uc.EmailService.SendOTP(ctx, req.Email, otp, RegistrationOTPExpiryMinutes)

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

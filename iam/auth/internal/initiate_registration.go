package internal

import (
	"context"
	"time"

	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

// InitiateRegistration creates a registration session in Redis and sends OTP to email.
// Design reference: .claude/doc/email-otp-signup-api.md Section 2.3
func (uc *usecase) InitiateRegistration(
	ctx context.Context,
	tenantID uuid.UUID,
	req *authdto.InitiateRegistrationRequest,
	ipAddress, userAgent string,
) (*authdto.InitiateRegistrationResponse, error) {
	// 1. Verify tenant exists
	tenantExists, err := uc.TenantRepo.Exists(ctx, tenantID)
	if err != nil {
		return nil, errors.ErrInternal("failed to verify tenant").WithError(err)
	}
	if !tenantExists {
		return nil, errors.ErrTenantNotFound()
	}

	// 2. Check if email already exists in this tenant
	emailExists, err := uc.UserRepo.EmailExistsInTenant(ctx, tenantID, req.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to check email").WithError(err)
	}
	if emailExists {
		return nil, errors.ErrUserAlreadyExists()
	}

	// 3. Check rate limit (returns AppError directly)
	rateLimitTTL := time.Duration(RegistrationRateLimitWindow) * time.Minute
	count, err := uc.Redis.IncrementRegistrationRateLimit(ctx, tenantID, req.Email, rateLimitTTL)
	if err != nil {
		return nil, err
	}
	if count > int64(RegistrationRateLimitPerHour) {
		return nil, errors.ErrTooManyRequests("Too many registration attempts. Please try again later.")
	}

	// 4. Check if there's already an active registration for this email (returns AppError directly)
	emailLocked, err := uc.Redis.IsRegistrationEmailLocked(ctx, tenantID, req.Email)
	if err != nil {
		return nil, err
	}
	if emailLocked {
		return nil, errors.ErrConflict("An active registration already exists for this email")
	}

	// 5. Generate OTP
	otp, otpHash, err := uc.generateOTP()
	if err != nil {
		return nil, errors.ErrInternal("failed to generate OTP").WithError(err)
	}

	// 6. Create registration session
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

	// 7. Lock the email to prevent duplicate registrations
	locked, err := uc.Redis.LockRegistrationEmail(ctx, tenantID, req.Email, sessionTTL)
	if err != nil {
		return nil, errors.ErrInternal("failed to lock email").WithError(err)
	}
	if !locked {
		return nil, errors.ErrConflict("An active registration already exists for this email")
	}

	// 8. Store session in Redis (returns AppError directly)
	if err := uc.Redis.CreateRegistrationSession(ctx, session, sessionTTL); err != nil {
		// Cleanup email lock on failure
		_ = uc.Redis.UnlockRegistrationEmail(ctx, tenantID, req.Email)
		return nil, err
	}

	// 9. Send OTP email
	if err := uc.EmailService.SendOTP(ctx, req.Email, otp, RegistrationOTPExpiryMinutes); err != nil {
		// Log error but don't fail - session is created
		// In production, you might want to implement retry logic
	}

	// 10. Return response
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

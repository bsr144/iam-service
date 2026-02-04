package internal

import (
	"context"
	"net/http"
	"time"

	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) ResendRegistrationOTP(
	ctx context.Context,
	tenantID, registrationID uuid.UUID,
	req *authdto.ResendRegistrationOTPRequest,
) (*authdto.ResendRegistrationOTPResponse, error) {
	session, err := uc.Redis.GetRegistrationSession(ctx, tenantID, registrationID)
	if err != nil {
		return nil, err
	}

	if session.Email != req.Email {
		return nil, errors.ErrValidation("Email does not match registration")
	}

	if session.IsExpired() {
		return nil, errors.New("REGISTRATION_EXPIRED", "Registration session has expired", http.StatusGone)
	}

	if session.Status == entity.RegistrationSessionStatusVerified {
		return nil, errors.ErrConflict("Registration is already verified")
	}

	if session.Status == entity.RegistrationSessionStatusFailed {
		return nil, errors.New("REGISTRATION_FAILED", "Registration has failed due to too many attempts. Please start a new registration.", http.StatusGone)
	}

	if session.Status != entity.RegistrationSessionStatusPendingVerification {
		return nil, errors.ErrBadRequest("Registration is not in a state where OTP can be resent")
	}

	if session.ResendCount >= session.MaxResends {
		return nil, errors.ErrTooManyRequests("Maximum number of resends reached")
	}

	if !session.CanResendOTP() {
		cooldownRemaining := session.CooldownRemainingSeconds()
		retryAfter := time.Now().Add(time.Duration(cooldownRemaining) * time.Second)
		return nil, errors.ErrTooManyRequests("Please wait before requesting another code").
			WithDetails(map[string]interface{}{
				"retry_after_seconds": cooldownRemaining,
				"retry_after":         retryAfter.Format(time.RFC3339),
			})
	}

	otp, otpHash, err := uc.generateOTP()
	if err != nil {
		return nil, errors.ErrInternal("failed to generate OTP").WithError(err)
	}

	otpExpiry := time.Now().Add(time.Duration(RegistrationOTPExpiryMinutes) * time.Minute)
	if err := uc.Redis.UpdateRegistrationOTP(ctx, tenantID, registrationID, otpHash, otpExpiry); err != nil {
		return nil, err
	}

	if err := uc.EmailService.SendOTP(ctx, req.Email, otp, RegistrationOTPExpiryMinutes); err != nil {

	}

	nextResendAt := time.Now().Add(time.Duration(session.ResendCooldownSeconds) * time.Second)

	return &authdto.ResendRegistrationOTPResponse{
		RegistrationID:        registrationID.String(),
		Message:               "New verification code sent to your email",
		ExpiresAt:             otpExpiry,
		ResendsRemaining:      session.MaxResends - session.ResendCount - 1,
		NextResendAvailableAt: nextResendAt,
	}, nil
}

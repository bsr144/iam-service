package internal

import (
	"context"
	"strings"

	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

// GetRegistrationStatus returns the current status of a registration session.
// Design reference: .claude/doc/email-otp-signup-api.md Section 2.7
func (uc *usecase) GetRegistrationStatus(
	ctx context.Context,
	tenantID, registrationID uuid.UUID,
	email string,
) (*authdto.RegistrationStatusResponse, error) {
	// 1. Get registration session (returns AppError directly)
	session, err := uc.Redis.GetRegistrationSession(ctx, tenantID, registrationID)
	if err != nil {
		return nil, err
	}

	// 2. Validate email matches (basic security check)
	if !strings.EqualFold(session.Email, email) {
		return nil, errors.ErrNotFound("Registration session not found")
	}

	// 3. Mask email for privacy (u***@example.com)
	maskedEmail := maskEmailForRegistration(session.Email)

	// 4. Return response
	return &authdto.RegistrationStatusResponse{
		RegistrationID:       registrationID.String(),
		Email:                maskedEmail,
		Status:               string(session.Status),
		ExpiresAt:            session.ExpiresAt,
		OTPAttemptsRemaining: session.RemainingAttempts(),
		ResendsRemaining:     session.RemainingResends(),
	}, nil
}

// maskEmailForRegistration masks an email address for privacy in registration context.
// Example: "user@example.com" -> "u***@example.com"
func maskEmailForRegistration(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***"
	}

	local := parts[0]
	domain := parts[1]

	if len(local) == 0 {
		return "***@" + domain
	}

	if len(local) == 1 {
		return local + "***@" + domain
	}

	return string(local[0]) + "***@" + domain
}

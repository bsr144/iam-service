package internal

import (
	"context"
	"strings"

	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) GetRegistrationStatus(
	ctx context.Context,
	tenantID, registrationID uuid.UUID,
	email string,
) (*authdto.RegistrationStatusResponse, error) {
	session, err := uc.Redis.GetRegistrationSession(ctx, tenantID, registrationID)
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(session.Email, email) {
		return nil, errors.ErrNotFound("Registration session not found")
	}

	maskedEmail := maskEmailForRegistration(session.Email)

	return &authdto.RegistrationStatusResponse{
		RegistrationID:       registrationID.String(),
		Email:                maskedEmail,
		Status:               string(session.Status),
		ExpiresAt:            session.ExpiresAt,
		OTPAttemptsRemaining: session.RemainingAttempts(),
		ResendsRemaining:     session.RemainingResends(),
	}, nil
}

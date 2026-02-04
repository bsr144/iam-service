package internal

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// VerifyRegistrationOTP verifies the OTP code and returns a registration token.
// Design reference: .claude/doc/email-otp-signup-api.md Section 2.4
func (uc *usecase) VerifyRegistrationOTP(
	ctx context.Context,
	tenantID, registrationID uuid.UUID,
	req *authdto.VerifyRegistrationOTPRequest,
) (*authdto.VerifyRegistrationOTPResponse, error) {
	// 1. Get registration session (returns AppError directly)
	session, err := uc.Redis.GetRegistrationSession(ctx, tenantID, registrationID)
	if err != nil {
		return nil, err
	}

	// 2. Validate email matches
	if session.Email != req.Email {
		return nil, errors.ErrValidation("Email does not match registration")
	}

	// 3. Check session status
	if session.IsExpired() {
		return nil, errors.New("REGISTRATION_EXPIRED", "Registration session has expired", http.StatusGone)
	}

	if session.Status == entity.RegistrationSessionStatusVerified {
		return nil, errors.ErrConflict("Registration is already verified")
	}

	if session.Status == entity.RegistrationSessionStatusFailed {
		return nil, errors.ErrTooManyRequests("Too many failed attempts. Please start a new registration.")
	}

	if session.Status != entity.RegistrationSessionStatusPendingVerification {
		return nil, errors.ErrBadRequest("Registration is not in a verifiable state")
	}

	// 4. Check if OTP has expired
	if session.IsOTPExpired() {
		return nil, errors.New("OTP_EXPIRED", "Verification code has expired. Please request a new one.", http.StatusGone)
	}

	// 5. Check attempt limit
	if !session.CanAttemptOTP() {
		return nil, errors.ErrTooManyRequests("Too many failed attempts. Please start a new registration.")
	}

	// 6. Verify OTP using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(session.OTPHash), []byte(req.OTPCode))
	if err != nil {
		// Increment attempts (returns AppError directly)
		attempts, incErr := uc.Redis.IncrementRegistrationAttempts(ctx, tenantID, registrationID)
		if incErr != nil {
			return nil, incErr
		}

		remaining := session.MaxAttempts - attempts
		if remaining <= 0 {
			return nil, errors.ErrTooManyRequests("Too many failed attempts. Registration has been invalidated.")
		}

		return nil, errors.ErrUnauthorized("The verification code is incorrect").
			WithDetails(map[string]interface{}{
				"attempts_remaining": remaining,
			})
	}

	// 7. Generate registration completion token
	token, tokenHash, err := uc.generateRegistrationCompleteToken(registrationID, tenantID, req.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate registration token").WithError(err)
	}

	// 8. Mark session as verified (returns AppError directly)
	if err := uc.Redis.MarkRegistrationVerified(ctx, tenantID, registrationID, tokenHash); err != nil {
		return nil, err
	}

	// 9. Return response
	tokenExpiry := time.Now().Add(time.Duration(RegistrationCompleteTokenExpiryMinutes) * time.Minute)

	return &authdto.VerifyRegistrationOTPResponse{
		RegistrationID:    registrationID.String(),
		Status:            string(entity.RegistrationSessionStatusVerified),
		Message:           "Email verified successfully",
		RegistrationToken: token,
		TokenExpiresAt:    tokenExpiry,
		NextStep: authdto.NextStep{
			Action:         "complete_registration",
			Endpoint:       fmt.Sprintf("/api/iam/v1/registrations/%s/complete", registrationID.String()),
			RequiredFields: []string{"password", "password_confirmation", "first_name", "last_name"},
		},
	}, nil
}

// generateRegistrationCompleteToken creates a JWT token for completing registration.
// Returns the token string and its SHA256 hash (for storage and single-use verification).
func (uc *usecase) generateRegistrationCompleteToken(registrationID, tenantID uuid.UUID, email string) (string, string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"registration_id": registrationID.String(),
		"tenant_id":       tenantID.String(),
		"email":           email,
		"purpose":         RegistrationCompleteTokenPurpose,
		"exp":             now.Add(time.Duration(RegistrationCompleteTokenExpiryMinutes) * time.Minute).Unix(),
		"iat":             now.Unix(),
		"jti":             uuid.New().String(), // Unique token ID for single-use
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(uc.Config.JWT.AccessSecret))
	if err != nil {
		return "", "", err
	}

	// Generate hash for storage (single-use verification)
	hash := sha256.Sum256([]byte(tokenString))
	tokenHash := hex.EncodeToString(hash[:])

	return tokenString, tokenHash, nil
}

// validateRegistrationCompleteToken validates the registration completion token.
// Returns the claims if valid.
func (uc *usecase) validateRegistrationCompleteToken(tokenString string, expectedRegistrationID, expectedTenantID uuid.UUID) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrTokenInvalid()
		}
		return []byte(uc.Config.JWT.AccessSecret), nil
	})

	if err != nil {
		return nil, errors.ErrUnauthorized("Registration token is invalid or expired")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.ErrUnauthorized("Registration token is invalid")
	}

	// Validate purpose
	if purpose, ok := claims["purpose"].(string); !ok || purpose != RegistrationCompleteTokenPurpose {
		return nil, errors.ErrUnauthorized("Token is not a registration completion token")
	}

	// Validate registration_id matches
	if regID, ok := claims["registration_id"].(string); !ok || regID != expectedRegistrationID.String() {
		return nil, errors.ErrUnauthorized("Token does not match this registration")
	}

	// Validate tenant_id matches
	if tid, ok := claims["tenant_id"].(string); !ok || tid != expectedTenantID.String() {
		return nil, errors.ErrUnauthorized("Token does not match this tenant")
	}

	return claims, nil
}

package response

import (
	"time"

	"github.com/google/uuid"
)

type LoginResponse struct {
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	ExpiresIn    int               `json:"expires_in"`
	TokenType    string            `json:"token_type"`
	User         LoginUserResponse `json:"user"`
}

type LoginUserResponse struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	Roles    []string  `json:"roles"`
}

type RegisterResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	Email        string    `json:"email"`
	Status       string    `json:"status"`
	OTPExpiresAt time.Time `json:"otp_expires_at"`
}

type RegisterSpecialAccountResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}

type VerifyOTPResponse struct {
	RegistrationToken string `json:"registration_token"`
	ExpiresIn         int    `json:"expires_in"`
	NextStep          string `json:"next_step"`
}

type CompleteProfileResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Status   string    `json:"status"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	Message  string    `json:"message"`
}

type ResendOTPResponse struct {
	OTPExpiresAt time.Time `json:"otp_expires_at"`
}

type SetupPINResponse struct {
	PINSetAt time.Time `json:"pin_set_at"`
}

type RequestPasswordResetResponse struct {
	OTPExpiresAt time.Time `json:"otp_expires_at"`
	EmailMasked  string    `json:"email_masked"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

type InitiateRegistrationOTPConfig struct {
	ExpiresInMinutes      int `json:"expires_in_minutes"`
	ResendCooldownSeconds int `json:"resend_cooldown_seconds"`
}

type InitiateRegistrationResponse struct {
	RegistrationID string                        `json:"registration_id"`
	Email          string                        `json:"email"`
	Status         string                        `json:"status"`
	Message        string                        `json:"message"`
	ExpiresAt      time.Time                     `json:"expires_at"`
	OTPConfig      InitiateRegistrationOTPConfig `json:"otp_config"`
}

type VerifyRegistrationOTPNextStep struct {
	Action   string `json:"action"`
	Endpoint string `json:"endpoint"`
}

type VerifyRegistrationOTPResponse struct {
	RegistrationID    string                        `json:"registration_id"`
	Status            string                        `json:"status"`
	Message           string                        `json:"message"`
	RegistrationToken string                        `json:"registration_token"`
	TokenExpiresAt    time.Time                     `json:"token_expires_at"`
	NextStep          VerifyRegistrationOTPNextStep `json:"next_step"`
}

type ResendRegistrationOTPResponse struct {
	RegistrationID        string    `json:"registration_id"`
	Message               string    `json:"message"`
	ExpiresAt             time.Time `json:"expires_at"`
	ResendsRemaining      int       `json:"resends_remaining"`
	NextResendAvailableAt time.Time `json:"next_resend_available_at"`
}

type RegistrationStatusResponse struct {
	RegistrationID       string    `json:"registration_id"`
	Email                string    `json:"email"`
	Status               string    `json:"status"`
	ExpiresAt            time.Time `json:"expires_at"`
	OTPAttemptsRemaining int       `json:"otp_attempts_remaining"`
	ResendsRemaining     int       `json:"resends_remaining"`
}

type CompleteRegistrationProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CompleteRegistrationResponse struct {
	UserID  uuid.UUID                   `json:"user_id"`
	Email   string                      `json:"email"`
	Status  string                      `json:"status"`
	Message string                      `json:"message"`
	Profile CompleteRegistrationProfile `json:"profile"`

	AccessToken  *string `json:"access_token,omitempty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
	TokenType    *string `json:"token_type,omitempty"`
	ExpiresIn    *int    `json:"expires_in,omitempty"`
}

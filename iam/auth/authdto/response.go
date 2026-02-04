package authdto

import (
	"time"

	"github.com/google/uuid"
)

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

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"`
	TokenType    string       `json:"token_type"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID         uuid.UUID  `json:"id"`
	Email      string     `json:"email"`
	FullName   string     `json:"full_name"`
	TenantID   *uuid.UUID `json:"tenant_id,omitempty"`
	ProductID  *uuid.UUID `json:"product_id,omitempty"`
	BranchID   *uuid.UUID `json:"branch_id,omitempty"`
	Roles      []string   `json:"roles"`
	MFAEnabled bool       `json:"mfa_enabled"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type SetupPINResponse struct {
	PINSetAt     time.Time `json:"pin_set_at"`
	PINExpiresAt time.Time `json:"pin_expires_at"`
}

type VerifyPINResponse struct {
	PINToken  string `json:"pin_token"`
	ExpiresIn int    `json:"expires_in"`
	Operation string `json:"operation"`
}

type ResendOTPResponse struct {
	OTPExpiresAt time.Time `json:"otp_expires_at"`
}

type RequestPINResetOTPResponse struct {
	OTPExpiresAt time.Time `json:"otp_expires_at"`
	EmailMasked  string    `json:"email_masked"`
}

type ResetPINResponse struct {
	PINChangedAt time.Time `json:"pin_changed_at"`
	PINExpiresAt time.Time `json:"pin_expires_at"`
}

type RequestPasswordResetResponse struct {
	OTPExpiresAt time.Time `json:"otp_expires_at"`
	EmailMasked  string    `json:"email_masked"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

type OTPConfig struct {
	Length                int `json:"length"`
	ExpiresInMinutes      int `json:"expires_in_minutes"`
	MaxAttempts           int `json:"max_attempts"`
	ResendCooldownSeconds int `json:"resend_cooldown_seconds"`
	MaxResends            int `json:"max_resends"`
}

type InitiateRegistrationResponse struct {
	RegistrationID string    `json:"registration_id"`
	Email          string    `json:"email"`
	Status         string    `json:"status"`
	Message        string    `json:"message"`
	ExpiresAt      time.Time `json:"expires_at"`
	OTPConfig      OTPConfig `json:"otp_config"`
}

type NextStep struct {
	Action         string   `json:"action"`
	Endpoint       string   `json:"endpoint"`
	RequiredFields []string `json:"required_fields"`
}

type VerifyRegistrationOTPResponse struct {
	RegistrationID    string    `json:"registration_id"`
	Status            string    `json:"status"`
	Message           string    `json:"message"`
	RegistrationToken string    `json:"registration_token"`
	TokenExpiresAt    time.Time `json:"token_expires_at"`
	NextStep          NextStep  `json:"next_step"`
}

type ResendRegistrationOTPResponse struct {
	RegistrationID        string    `json:"registration_id"`
	Message               string    `json:"message"`
	ExpiresAt             time.Time `json:"expires_at"`
	ResendsRemaining      int       `json:"resends_remaining"`
	NextResendAvailableAt time.Time `json:"next_resend_available_at"`
}

type RegistrationUserProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CompleteRegistrationResponse struct {
	UserID  uuid.UUID               `json:"user_id"`
	Email   string                  `json:"email"`
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Profile RegistrationUserProfile `json:"profile"`

	AccessToken  *string `json:"access_token,omitempty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
	TokenType    *string `json:"token_type,omitempty"`
	ExpiresIn    *int    `json:"expires_in,omitempty"`
}

type RegistrationStatusResponse struct {
	RegistrationID       string    `json:"registration_id"`
	Email                string    `json:"email"`
	Status               string    `json:"status"`
	ExpiresAt            time.Time `json:"expires_at"`
	OTPAttemptsRemaining int       `json:"otp_attempts_remaining"`
	ResendsRemaining     int       `json:"resends_remaining"`
}

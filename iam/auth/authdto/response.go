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

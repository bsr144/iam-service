package authdto

import (
	"github.com/google/uuid"
)

type RegisterRequest struct {
	TenantID  uuid.UUID `json:"tenant_id" validate:"required"`
	Email     string    `json:"email" validate:"required,email,max=255"`
	Password  string    `json:"password" validate:"required,min=8,max=128"`
	FirstName string    `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string    `json:"last_name" validate:"required,min=2,max=100"`
}
type RegisterSpecialAccountRequest struct {
	TenantID  uuid.UUID `json:"tenant_id" validate:"required"`
	UserType  string    `json:"user_type" validate:"required,oneof=ADMIN APPROVER"`
	Email     string    `json:"email" validate:"required,email,max=255"`
	Password  string    `json:"password" validate:"required,min=8,max=128"`
	FirstName string    `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string    `json:"last_name" validate:"required,min=2,max=100"`
}

type VerifyOTPRequest struct {
	TenantID uuid.UUID `json:"tenant_id" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
	OTPCode  string    `json:"otp_code" validate:"required,len=6,numeric"`
}

type ResendOTPRequest struct {
	TenantID uuid.UUID `json:"tenant_id" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
}

type CompleteProfileRequest struct {
	RegistrationToken string  `json:"registration_token" validate:"required"`
	Address           *string `json:"address,omitempty" validate:"omitempty,min=10,max=500"`
	Phone             *string `json:"phone,omitempty" validate:"omitempty"`
	Gender            *string `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	MaritalStatus     *string `json:"marital_status,omitempty" validate:"omitempty,oneof=single married divorced widowed"`
	DateOfBirth       *string `json:"date_of_birth,omitempty"`
	PlaceOfBirth      *string `json:"place_of_birth,omitempty" validate:"omitempty,min=2,max=100"`
}

type LoginRequest struct {
	TenantID    uuid.UUID  `json:"tenant_id" validate:"required"`
	Email       string     `json:"email" validate:"required,email"`
	Password    string     `json:"password" validate:"required"`
	ProductID   *uuid.UUID `json:"product_id,omitempty"`
	ProductCode *string    `json:"product_code,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type SetupPINRequest struct {
	PIN        string `json:"pin" validate:"required,len=6,numeric"`
	PINConfirm string `json:"pin_confirm" validate:"required,len=6,numeric"`
}

type VerifyPINRequest struct {
	PIN       string `json:"pin" validate:"required,len=6,numeric"`
	Operation string `json:"operation" validate:"required"`
}

type RequestPINResetOTPRequest struct {
	Password string `json:"password" validate:"required"`
}

type ResetPINRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	OTP             string `json:"otp" validate:"required,len=6,numeric"`
	NewPIN          string `json:"new_pin" validate:"required,len=6,numeric"`
	NewPINConfirm   string `json:"new_pin_confirm" validate:"required,len=6,numeric"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RequestPasswordResetRequest struct {
	TenantID uuid.UUID `json:"tenant_id" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	TenantID    uuid.UUID `json:"tenant_id" validate:"required"`
	Email       string    `json:"email" validate:"required,email"`
	OTPCode     string    `json:"otp_code" validate:"required,len=6,numeric"`
	NewPassword string    `json:"new_password" validate:"required,min=8,max=128"`
}

type InitiateRegistrationRequest struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type VerifyRegistrationOTPRequest struct {
	Email   string `json:"email" validate:"required,email"`
	OTPCode string `json:"otp_code" validate:"required,len=6,numeric"`
}

type ResendRegistrationOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type CompleteRegistrationRequest struct {
	Password             string  `json:"password" validate:"required,min=8,max=128"`
	PasswordConfirmation string  `json:"password_confirmation" validate:"required,eqfield=Password"`
	FirstName            string  `json:"first_name" validate:"required,min=1,max=100"`
	LastName             string  `json:"last_name" validate:"required,min=1,max=100"`
	PhoneNumber          *string `json:"phone_number,omitempty" validate:"omitempty,e164"`
}

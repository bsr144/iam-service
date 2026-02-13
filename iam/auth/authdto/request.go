package authdto

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
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

type SetPasswordRequest struct {
	Password             string `json:"password" validate:"required,min=8,max=128"`
	ConfirmationPassword string `json:"confirmation_password" validate:"required,eqfield=Password"`
}

type CompleteProfileRegistrationRequest struct {
	FullName      string `json:"full_name" validate:"required,min=1,max=200"`
	PhoneNumber   string `json:"phone_number" validate:"required,e164"`
	DateOfBirth   string `json:"date_of_birth" validate:"required"`
	Gender        string `json:"gender" validate:"required,oneof=male female other"`
	MaritalStatus string `json:"marital_status" validate:"required,oneof=single married divorced widowed"`
	Address       string `json:"address" validate:"required,min=10,max=500"`
	PlaceOfBirth  string `json:"place_of_birth" validate:"required,min=2,max=100"`
}

type CompleteRegistrationRequest struct {
	Password             string  `json:"password" validate:"required,min=8,max=128"`
	PasswordConfirmation string  `json:"password_confirmation" validate:"required,eqfield=Password"`
	FirstName            string  `json:"first_name" validate:"required,min=1,max=100"`
	LastName             string  `json:"last_name" validate:"required,min=1,max=100"`
	PhoneNumber          *string `json:"phone_number,omitempty" validate:"omitempty,e164"`
}

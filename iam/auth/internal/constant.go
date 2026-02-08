package internal

const (
	PasswordMinLength = 8
)

const (
	OTPLength        = 6
	OTPExpiryMinutes = 1

	MaxActiveOTPPerEmail = 5
)

const (
	JWTExpirationHours = 72
)

const (
	RegistrationTokenExpiryMinutes = 10
)

const (
	RegistrationSessionExpiryMinutes = 10

	RegistrationOTPLength         = 6
	RegistrationOTPExpiryMinutes  = 10
	RegistrationOTPMaxAttempts    = 5
	RegistrationOTPMaxResends     = 3
	RegistrationOTPResendCooldown = 60

	RegistrationCompleteTokenExpiryMinutes = 15
	RegistrationCompleteTokenPurpose       = "registration_complete"

	RegistrationRateLimitPerHour = 3
	RegistrationRateLimitWindow  = 60
)

const (
	LoginSessionExpiryMinutes = 10

	LoginOTPLength         = 6
	LoginOTPExpiryMinutes  = 10
	LoginOTPMaxAttempts    = 5
	LoginOTPMaxResends     = 3
	LoginOTPResendCooldown = 60

	LoginRateLimitPerHour = 10
)

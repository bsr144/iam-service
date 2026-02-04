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

// =============================================================================
// Email + OTP Registration Flow (Redis-based)
// Design reference: .claude/doc/email-otp-signup-api.md Section 6.1
// =============================================================================

const (
	// Registration session settings
	RegistrationSessionExpiryMinutes = 10

	// Registration OTP settings
	RegistrationOTPLength          = 6
	RegistrationOTPExpiryMinutes   = 10
	RegistrationOTPMaxAttempts     = 5
	RegistrationOTPMaxResends      = 3
	RegistrationOTPResendCooldown  = 60 // seconds

	// Registration token settings (for completing registration after OTP verification)
	RegistrationCompleteTokenExpiryMinutes = 15
	RegistrationCompleteTokenPurpose       = "registration_complete"

	// Rate limiting
	RegistrationRateLimitPerHour = 3
	RegistrationRateLimitWindow  = 60 // minutes
)

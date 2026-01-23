package entity

import (
	"encoding/json"
	"net"
	"time"

	"github.com/google/uuid"
)

// VerificationEntityType represents what type of entity is being verified
type VerificationEntityType string

const (
	VerificationEntityTypeRegistration  VerificationEntityType = "registration"
	VerificationEntityTypeUser          VerificationEntityType = "user"
	VerificationEntityTypePasswordReset VerificationEntityType = "password_reset"
	VerificationEntityTypeEmailChange   VerificationEntityType = "email_change"
	VerificationEntityTypeStepUpAuth    VerificationEntityType = "step_up_auth"
	VerificationEntityTypeSession       VerificationEntityType = "session"
)

// VerificationPurpose represents why verification is needed
type VerificationPurpose string

const (
	VerificationPurposeRegister           VerificationPurpose = "register"
	VerificationPurposeResetPassword      VerificationPurpose = "reset_password"
	VerificationPurposeChangeEmail        VerificationPurpose = "change_email"
	VerificationPurposeSensitiveOperation VerificationPurpose = "sensitive_operation"
	VerificationPurposeLoginMFA           VerificationPurpose = "login_mfa"
	VerificationPurposeStepUp             VerificationPurpose = "step_up"
)

// VerificationMethod represents how to verify
type VerificationMethod string

const (
	VerificationMethodOTPEmail  VerificationMethod = "otp_email"
	VerificationMethodOTPSMS    VerificationMethod = "otp_sms"
	VerificationMethodPIN       VerificationMethod = "pin"
	VerificationMethodTOTP      VerificationMethod = "totp"
	VerificationMethodBiometric VerificationMethod = "biometric"
	VerificationMethodLiveness  VerificationMethod = "liveness"
	VerificationMethodWebAuthn  VerificationMethod = "webauthn"
)

// VerificationStatus represents the state of a verification
type VerificationStatus string

const (
	VerificationStatusPending  VerificationStatus = "pending"
	VerificationStatusSent     VerificationStatus = "sent"
	VerificationStatusVerified VerificationStatus = "verified"
	VerificationStatusFailed   VerificationStatus = "failed"
	VerificationStatusExpired  VerificationStatus = "expired"
	VerificationStatusLocked   VerificationStatus = "locked"
)

// VerificationDeliveryStatus represents the delivery state
type VerificationDeliveryStatus string

const (
	VerificationDeliveryStatusPending   VerificationDeliveryStatus = "pending"
	VerificationDeliveryStatusSent      VerificationDeliveryStatus = "sent"
	VerificationDeliveryStatusFailed    VerificationDeliveryStatus = "failed"
	VerificationDeliveryStatusDelivered VerificationDeliveryStatus = "delivered"
	VerificationDeliveryStatusBounced   VerificationDeliveryStatus = "bounced"
)

// VerificationDeliveryChannel represents how the challenge is delivered
type VerificationDeliveryChannel string

const (
	VerificationDeliveryChannelEmail  VerificationDeliveryChannel = "email"
	VerificationDeliveryChannelSMS    VerificationDeliveryChannel = "sms"
	VerificationDeliveryChannelPush   VerificationDeliveryChannel = "push"
	VerificationDeliveryChannelApp    VerificationDeliveryChannel = "app"
	VerificationDeliveryChannelDevice VerificationDeliveryChannel = "device"
)

// Verification represents a generic verification framework
// Supports multiple verification methods: OTP, PIN, biometric, liveness, etc.
// Reusable across multiple purposes: registration, password reset, step-up auth, etc.
type Verification struct {
	VerificationID uuid.UUID `json:"verification_id" gorm:"column:verification_id;primaryKey" db:"verification_id"`

	// What is being verified? (polymorphic reference)
	EntityType VerificationEntityType `json:"entity_type" gorm:"column:entity_type;not null" db:"entity_type"`
	EntityID   uuid.UUID              `json:"entity_id" gorm:"column:entity_id;not null" db:"entity_id"`

	// Why verification is needed
	Purpose VerificationPurpose `json:"purpose" gorm:"column:purpose;not null" db:"purpose"`

	// How to verify (method agnostic)
	VerificationMethod VerificationMethod `json:"verification_method" gorm:"column:verification_method;not null" db:"verification_method"`

	// Delivery information
	DeliveryTarget          string                          `json:"delivery_target" gorm:"column:delivery_target;not null" db:"delivery_target"`
	DeliveryChannel         *VerificationDeliveryChannel    `json:"delivery_channel,omitempty" gorm:"column:delivery_channel" db:"delivery_channel"`
	DeliveryStatus          VerificationDeliveryStatus      `json:"delivery_status" gorm:"column:delivery_status;default:'pending'" db:"delivery_status"`
	DeliveryAttempts        int                             `json:"delivery_attempts" gorm:"column:delivery_attempts;default:0" db:"delivery_attempts"`
	LastDeliveryAttemptAt   *time.Time                      `json:"last_delivery_attempt_at,omitempty" gorm:"column:last_delivery_attempt_at" db:"last_delivery_attempt_at"`
	DeliveryError           *string                         `json:"delivery_error,omitempty" gorm:"column:delivery_error" db:"delivery_error"`

	// Security constraints
	MaxAttempts  int        `json:"max_attempts" gorm:"column:max_attempts;not null;default:3" db:"max_attempts"`
	AttemptsUsed int        `json:"attempts_used" gorm:"column:attempts_used;not null;default:0" db:"attempts_used"`
	LockedUntil  *time.Time `json:"locked_until,omitempty" gorm:"column:locked_until" db:"locked_until"`

	// State management
	Status VerificationStatus `json:"status" gorm:"column:status;not null;default:'pending'" db:"status"`

	// Context tracking
	IPAddress *net.IP         `json:"ip_address,omitempty" gorm:"column:ip_address;type:inet" db:"ip_address"`
	UserAgent *string         `json:"user_agent,omitempty" gorm:"column:user_agent" db:"user_agent"`
	Metadata  json.RawMessage `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb;default:'{}'" db:"metadata"`

	// Lifecycle timestamps
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" db:"created_at"`
	ExpiresAt  time.Time  `json:"expires_at" gorm:"column:expires_at;not null" db:"expires_at"`
	VerifiedAt *time.Time `json:"verified_at,omitempty" gorm:"column:verified_at" db:"verified_at"`

	// Result tracking
	VerificationResult json.RawMessage `json:"verification_result,omitempty" gorm:"column:verification_result;type:jsonb" db:"verification_result"`
	FailureReason      *string         `json:"failure_reason,omitempty" gorm:"column:failure_reason" db:"failure_reason"`
}

func (Verification) TableName() string {
	return "verifications"
}

// IsExpired checks if the verification has expired
func (v *Verification) IsExpired() bool {
	return time.Now().After(v.ExpiresAt)
}

// IsLocked checks if the verification is locked due to too many failed attempts
func (v *Verification) IsLocked() bool {
	if v.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*v.LockedUntil)
}

// CanAttempt checks if another verification attempt is allowed
func (v *Verification) CanAttempt() bool {
	return !v.IsExpired() && !v.IsLocked() && v.AttemptsUsed < v.MaxAttempts && v.Status == VerificationStatusSent
}

// IsVerified checks if verification has been completed successfully
func (v *Verification) IsVerified() bool {
	return v.Status == VerificationStatusVerified
}

// RemainingAttempts returns the number of attempts left
func (v *Verification) RemainingAttempts() int {
	remaining := v.MaxAttempts - v.AttemptsUsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// ChallengeType represents the type of verification challenge
type ChallengeType string

const (
	ChallengeTypeOTP              ChallengeType = "otp"
	ChallengeTypeTOTP             ChallengeType = "totp"
	ChallengeTypePIN              ChallengeType = "pin"
	ChallengeTypeBiometricHash    ChallengeType = "biometric_hash"
	ChallengeTypeLivenessToken    ChallengeType = "liveness_token"
	ChallengeTypeWebAuthnChallenge ChallengeType = "webauthn_challenge"
)

// VerificationChallenge stores the actual verification data (OTP code, PIN hash, etc.)
// Separates verification intent (Verification) from verification proof (VerificationChallenge)
type VerificationChallenge struct {
	ChallengeID    uuid.UUID `json:"challenge_id" gorm:"column:challenge_id;primaryKey" db:"challenge_id"`
	VerificationID uuid.UUID `json:"verification_id" gorm:"column:verification_id;not null" db:"verification_id"`

	// Challenge data (method-specific)
	ChallengeType ChallengeType   `json:"challenge_type" gorm:"column:challenge_type;not null" db:"challenge_type"`
	ChallengeHash *string         `json:"-" gorm:"column:challenge_hash" db:"challenge_hash"` // bcrypt hash for OTP/PIN
	ChallengeData json.RawMessage `json:"challenge_data,omitempty" gorm:"column:challenge_data;type:jsonb" db:"challenge_data"` // For complex challenges

	// OTP-specific fields (for logging and debugging)
	OTPCodePrefix *string `json:"otp_code_prefix,omitempty" gorm:"column:otp_code_prefix" db:"otp_code_prefix"` // First 2 digits
	OTPDeliveryID *string `json:"otp_delivery_id,omitempty" gorm:"column:otp_delivery_id" db:"otp_delivery_id"` // External service ID

	// Lifecycle timestamps
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" db:"created_at"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"column:expires_at;not null" db:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty" gorm:"column:used_at" db:"used_at"`

	// Security flags
	IsUsed bool `json:"is_used" gorm:"column:is_used;not null;default:false" db:"is_used"`
}

func (VerificationChallenge) TableName() string {
	return "verification_challenges"
}

// IsExpired checks if the challenge has expired
func (vc *VerificationChallenge) IsExpired() bool {
	return time.Now().After(vc.ExpiresAt)
}

// CanBeUsed checks if the challenge can be used for verification
func (vc *VerificationChallenge) CanBeUsed() bool {
	return !vc.IsUsed && !vc.IsExpired()
}

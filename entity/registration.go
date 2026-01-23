package entity

import (
	"encoding/json"
	"net"
	"time"

	"github.com/google/uuid"
)

// RegistrationStatus represents the lifecycle state of a registration
type RegistrationStatus string

const (
	RegistrationStatusPendingVerification RegistrationStatus = "pending_verification"
	RegistrationStatusVerified            RegistrationStatus = "verified"
	RegistrationStatusCompleted           RegistrationStatus = "completed"
	RegistrationStatusExpired             RegistrationStatus = "expired"
	RegistrationStatusCancelled           RegistrationStatus = "cancelled"
)

// Registration represents a registration process in progress, exists before user creation
// This prevents email hijacking by not committing email to users table until OTP verification
type Registration struct {
	RegistrationID uuid.UUID  `json:"registration_id" gorm:"column:registration_id;primaryKey" db:"registration_id"`
	TenantID       uuid.UUID  `json:"tenant_id" gorm:"column:tenant_id;not null" db:"tenant_id"`
	BranchID       *uuid.UUID `json:"branch_id,omitempty" gorm:"column:branch_id" db:"branch_id"`

	// Identity intent (not yet committed to users table)
	Email        string `json:"email" gorm:"column:email;not null" db:"email"`
	PasswordHash string `json:"-" gorm:"column:password_hash;not null" db:"password_hash"` // Temporarily stored

	// Context tracking
	UserAgent *string         `json:"user_agent,omitempty" gorm:"column:user_agent" db:"user_agent"`
	IPAddress *net.IP         `json:"ip_address,omitempty" gorm:"column:ip_address;type:inet" db:"ip_address"`
	Referrer  *string         `json:"referrer,omitempty" gorm:"column:referrer" db:"referrer"`
	Metadata  json.RawMessage `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb;default:'{}'" db:"metadata"`

	// State management
	Status RegistrationStatus `json:"status" gorm:"column:status;not null;default:'pending_verification'" db:"status"`

	// Lifecycle timestamps
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" db:"created_at"`
	ExpiresAt   time.Time  `json:"expires_at" gorm:"column:expires_at;not null" db:"expires_at"`
	VerifiedAt  *time.Time `json:"verified_at,omitempty" gorm:"column:verified_at" db:"verified_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" gorm:"column:completed_at" db:"completed_at"`

	// Result tracking
	UserID              *uuid.UUID `json:"user_id,omitempty" gorm:"column:user_id" db:"user_id"`
	CancellationReason  *string    `json:"cancellation_reason,omitempty" gorm:"column:cancellation_reason" db:"cancellation_reason"`
}

func (Registration) TableName() string {
	return "registrations"
}

// IsExpired checks if the registration has expired
func (r *Registration) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// IsPending checks if the registration is still pending verification
func (r *Registration) IsPending() bool {
	return r.Status == RegistrationStatusPendingVerification
}

// IsVerified checks if OTP verification has been completed
func (r *Registration) IsVerified() bool {
	return r.Status == RegistrationStatusVerified
}

// IsCompleted checks if the registration process is completed (user created)
func (r *Registration) IsCompleted() bool {
	return r.Status == RegistrationStatusCompleted
}

// CanBeVerified checks if the registration can still accept OTP verification
func (r *Registration) CanBeVerified() bool {
	return r.Status == RegistrationStatusPendingVerification && !r.IsExpired()
}

// CanBeCompleted checks if the registration can be completed (create user)
func (r *Registration) CanBeCompleted() bool {
	return r.Status == RegistrationStatusVerified && !r.IsExpired()
}

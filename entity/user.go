package entity

import (
	"encoding/json"
	"net"
	"time"

	"github.com/google/uuid"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type MaritalStatus string

const (
	MaritalStatusSingle   MaritalStatus = "single"
	MaritalStatusMarried  MaritalStatus = "married"
	MaritalStatusDivorced MaritalStatus = "divorced"
	MaritalStatusWidowed  MaritalStatus = "widowed"
)

type UserStatus string

const (
	UserStatusPendingOTPVerification UserStatus = "pending_otp_verification"
	UserStatusPendingAdminApproval   UserStatus = "pending_admin_approval"
	UserStatusPendingUserCompletion  UserStatus = "pending_user_completion"
	UserStatusPendingPasswordChange  UserStatus = "pending_password_change"
	UserStatusPendingPINSetup        UserStatus = "pending_pin_setup"
	UserStatusActive                 UserStatus = "active"
	UserStatusSuspended              UserStatus = "suspended"
	UserStatusDeleted                UserStatus = "deleted"
)

type User struct {
	UserID           uuid.UUID  `json:"user_id" gorm:"column:user_id;primaryKey;type:uuid;default:uuidv7()" db:"user_id"`
	TenantID         *uuid.UUID `json:"tenant_id,omitempty" gorm:"column:tenant_id" db:"tenant_id"`
	BranchID         *uuid.UUID `json:"branch_id,omitempty" gorm:"column:branch_id" db:"branch_id"`
	Email            string     `json:"email" gorm:"column:email;not null" db:"email"`
	EmailVerified    bool       `json:"email_verified" gorm:"column:email_verified;default:false" db:"email_verified"`
	EmailVerifiedAt  *time.Time `json:"email_verified_at,omitempty" gorm:"column:email_verified_at" db:"email_verified_at"`
	IsServiceAccount bool       `json:"is_service_account" gorm:"column:is_service_account;default:false" db:"is_service_account"`
	IsActive         bool       `json:"is_active" gorm:"column:is_active;default:true" db:"is_active"`

	RegistrationID          *uuid.UUID `json:"registration_id,omitempty" gorm:"column:registration_id" db:"registration_id"`
	RegistrationCompletedAt *time.Time `json:"registration_completed_at,omitempty" gorm:"column:registration_completed_at" db:"registration_completed_at"`

	Timestamps
}

func (User) TableName() string {
	return "users"
}

func (u *User) IsPlatformAdmin() bool {
	return u.TenantID == nil
}

func (u *User) IsTenantUser() bool {
	return u.TenantID != nil
}

type UserCredentials struct {
	UserCredentialID  uuid.UUID       `json:"user_credential_id" gorm:"column:user_credential_id;primaryKey;type:uuid;default:uuidv7()" db:"user_credential_id"`
	UserID            uuid.UUID       `json:"user_id" gorm:"column:user_id;uniqueIndex;not null" db:"user_id"`
	PasswordHash      *string         `json:"-" gorm:"column:password_hash" db:"password_hash"`
	PasswordChangedAt *time.Time      `json:"password_changed_at,omitempty" gorm:"column:password_changed_at" db:"password_changed_at"`
	PasswordExpiresAt *time.Time      `json:"password_expires_at,omitempty" gorm:"column:password_expires_at" db:"password_expires_at"`
	PasswordHistory   json.RawMessage `json:"-" gorm:"column:password_history;type:jsonb;default:'[]'" db:"password_history"`
	PINHash           *string         `json:"-" gorm:"column:pin_hash" db:"pin_hash"`
	PINSetAt          *time.Time      `json:"pin_set_at,omitempty" gorm:"column:pin_set_at" db:"pin_set_at"`
	PINChangedAt      *time.Time      `json:"pin_changed_at,omitempty" gorm:"column:pin_changed_at" db:"pin_changed_at"`
	PINExpiresAt      *time.Time      `json:"pin_expires_at,omitempty" gorm:"column:pin_expires_at" db:"pin_expires_at"`
	PINHistory        json.RawMessage `json:"-" gorm:"column:pin_history;type:jsonb;default:'[]'" db:"pin_history"`
	SSOProvider       *string         `json:"sso_provider,omitempty" gorm:"column:sso_provider" db:"sso_provider"`
	SSOProviderID     *string         `json:"sso_provider_id,omitempty" gorm:"column:sso_provider_id" db:"sso_provider_id"`
	MFAEnabled        bool            `json:"mfa_enabled" gorm:"column:mfa_enabled;default:false" db:"mfa_enabled"`
	CreatedAt         time.Time       `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
}

func (UserCredentials) TableName() string {
	return "user_credentials"
}

type UserProfile struct {
	UserProfileID     uuid.UUID      `json:"user_profile_id" gorm:"column:user_profile_id;primaryKey;type:uuid;default:uuidv7()" db:"user_profile_id"`
	UserID            uuid.UUID      `json:"user_id" gorm:"column:user_id;uniqueIndex;not null" db:"user_id"`
	FirstName         string         `json:"first_name" gorm:"column:first_name;not null" db:"first_name"`
	LastName          string         `json:"last_name" gorm:"column:last_name;not null" db:"last_name"`
	Address           *string        `json:"address,omitempty" gorm:"column:address" db:"address"`
	Phone             *string        `json:"phone,omitempty" gorm:"column:phone" db:"phone"`
	Gender            *Gender        `json:"gender,omitempty" gorm:"column:gender" db:"gender"`
	MaritalStatus     *MaritalStatus `json:"marital_status,omitempty" gorm:"column:marital_status" db:"marital_status"`
	DateOfBirth       *string        `json:"date_of_birth,omitempty" gorm:"column:date_of_birth" db:"date_of_birth"`
	PlaceOfBirth      *string        `json:"place_of_birth,omitempty" gorm:"column:place_of_birth" db:"place_of_birth"`
	AvatarURL         *string        `json:"avatar_url,omitempty" gorm:"column:avatar_url" db:"avatar_url"`
	PreferredLanguage string         `json:"preferred_language" gorm:"column:preferred_language;default:en" db:"preferred_language"`
	Timezone          string         `json:"timezone" gorm:"column:timezone;default:Asia/Jakarta" db:"timezone"`
	CreatedAt         time.Time      `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}

func (u *UserProfile) FullName() string {
	return u.FirstName + " " + u.LastName
}

type UserSecurity struct {
	UserSecurityID      uuid.UUID       `json:"user_security_id" gorm:"column:user_security_id;primaryKey;type:uuid;default:uuidv7()" db:"user_security_id"`
	UserID              uuid.UUID       `json:"user_id" gorm:"column:user_id;uniqueIndex;not null" db:"user_id"`
	LastLoginAt         *time.Time      `json:"last_login_at,omitempty" gorm:"column:last_login_at" db:"last_login_at"`
	LastLoginIP         net.IP          `json:"last_login_ip,omitempty" gorm:"column:last_login_ip" db:"last_login_ip"`
	FailedLoginAttempts int             `json:"failed_login_attempts" gorm:"column:failed_login_attempts;default:0" db:"failed_login_attempts"`
	LockedUntil         *time.Time      `json:"locked_until,omitempty" gorm:"column:locked_until" db:"locked_until"`
	AdminRegisteredAt   *time.Time      `json:"admin_registered_at,omitempty" gorm:"column:admin_registered_at" db:"admin_registered_at"`
	UserRegisteredAt    *time.Time      `json:"user_registered_at,omitempty" gorm:"column:user_registered_at" db:"user_registered_at"`
	AdminRegisteredBy   *uuid.UUID      `json:"admin_registered_by,omitempty" gorm:"column:admin_registered_by" db:"admin_registered_by"`
	InvitationTokenHash *string         `json:"-" gorm:"column:invitation_token_hash" db:"invitation_token_hash"`
	InvitationExpiresAt *time.Time      `json:"invitation_expires_at,omitempty" gorm:"column:invitation_expires_at" db:"invitation_expires_at"`
	Metadata            json.RawMessage `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb;default:'{}'" db:"metadata"`
	CreatedAt           time.Time       `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
}

func (UserSecurity) TableName() string {
	return "user_security"
}

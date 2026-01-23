package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ActivationMethod string

const (
	ActivationMethodAdminFirst   ActivationMethod = "admin_first"
	ActivationMethodUserFirst    ActivationMethod = "user_first"
	ActivationMethodSimultaneous ActivationMethod = "simultaneous"
)

type StatusTransition struct {
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	TriggeredBy string    `json:"triggered_by"`
}

type UserActivationTracking struct {
	UserActivationTrackingID uuid.UUID  `json:"user_activation_tracking_id" gorm:"column:user_activation_tracking_id;primaryKey" db:"user_activation_tracking_id"`
	UserID                   uuid.UUID  `json:"user_id" gorm:"column:user_id;uniqueIndex;not null" db:"user_id"`
	TenantID                 *uuid.UUID `json:"tenant_id,omitempty" gorm:"column:tenant_id" db:"tenant_id"`

	AdminCreated   bool       `json:"admin_created" gorm:"column:admin_created;default:false" db:"admin_created"`
	AdminCreatedAt *time.Time `json:"admin_created_at,omitempty" gorm:"column:admin_created_at" db:"admin_created_at"`
	AdminCreatedBy *uuid.UUID `json:"admin_created_by,omitempty" gorm:"column:admin_created_by" db:"admin_created_by"`

	UserCompleted      bool       `json:"user_completed" gorm:"column:user_completed;default:false" db:"user_completed"`
	UserCompletedAt    *time.Time `json:"user_completed_at,omitempty" gorm:"column:user_completed_at" db:"user_completed_at"`
	OTPVerifiedAt      *time.Time `json:"otp_verified_at,omitempty" gorm:"column:otp_verified_at" db:"otp_verified_at"`
	ProfileCompletedAt *time.Time `json:"profile_completed_at,omitempty" gorm:"column:profile_completed_at" db:"profile_completed_at"`
	PINSetAt           *time.Time `json:"pin_set_at,omitempty" gorm:"column:pin_set_at" db:"pin_set_at"`

	ActivatedAt      *time.Time        `json:"activated_at,omitempty" gorm:"column:activated_at" db:"activated_at"`
	ActivationMethod *ActivationMethod `json:"activation_method,omitempty" gorm:"column:activation_method" db:"activation_method"`

	StatusHistory json.RawMessage `json:"status_history" gorm:"column:status_history;type:jsonb;default:'[]'" db:"status_history"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
}

func (UserActivationTracking) TableName() string {
	return "user_activation_tracking"
}

func (u *UserActivationTracking) IsActivated() bool {
	return u.ActivatedAt != nil
}

func (u *UserActivationTracking) IsAdminRegistered() bool {
	return u.AdminCreated && u.AdminCreatedAt != nil
}

func (u *UserActivationTracking) IsUserRegistered() bool {
	return u.UserCompleted && u.UserCompletedAt != nil
}

func (u *UserActivationTracking) IsPendingUserRegistration() bool {
	return u.IsAdminRegistered() && !u.IsUserRegistered()
}

func (u *UserActivationTracking) IsPendingAdminApproval() bool {
	return u.IsUserRegistered() && !u.IsAdminRegistered()
}

func (u *UserActivationTracking) GetStatusHistory() ([]StatusTransition, error) {
	var history []StatusTransition
	if u.StatusHistory == nil {
		return history, nil
	}
	if err := json.Unmarshal(u.StatusHistory, &history); err != nil {
		return nil, err
	}
	return history, nil
}

func (u *UserActivationTracking) AddStatusTransition(status, triggeredBy string) error {
	history, err := u.GetStatusHistory()
	if err != nil {
		history = []StatusTransition{}
	}

	transition := StatusTransition{
		Status:      status,
		Timestamp:   time.Now(),
		TriggeredBy: triggeredBy,
	}
	history = append(history, transition)

	data, err := json.Marshal(history)
	if err != nil {
		return err
	}
	u.StatusHistory = data
	u.UpdatedAt = time.Now()
	return nil
}

func NewUserActivationTracking(userID uuid.UUID, tenantID *uuid.UUID) *UserActivationTracking {
	now := time.Now()
	emptyHistory, _ := json.Marshal([]StatusTransition{})
	return &UserActivationTracking{
		UserActivationTrackingID: uuid.New(),
		UserID:                   userID,
		TenantID:                 tenantID,
		AdminCreated:             false,
		UserCompleted:            false,
		StatusHistory:            emptyHistory,
		CreatedAt:                now,
		UpdatedAt:                now,
	}
}

func (u *UserActivationTracking) MarkAdminCreated(adminID uuid.UUID) error {
	now := time.Now()
	u.AdminCreated = true
	u.AdminCreatedAt = &now
	u.AdminCreatedBy = &adminID
	return u.AddStatusTransition("admin_registered", "admin")
}
func (u *UserActivationTracking) MarkUserCreatedBySystem() error {
	now := time.Now()
	u.AdminCreated = true
	u.AdminCreatedAt = &now
	return u.AddStatusTransition("admin_registered", "system")
}

func (u *UserActivationTracking) MarkOTPVerified() error {
	now := time.Now()
	u.OTPVerifiedAt = &now
	return u.AddStatusTransition("otp_verified", "user")
}

func (u *UserActivationTracking) MarkProfileCompleted() error {
	now := time.Now()
	u.ProfileCompletedAt = &now
	return u.AddStatusTransition("profile_completed", "user")
}

func (u *UserActivationTracking) MarkPINSet() error {
	now := time.Now()
	u.PINSetAt = &now
	return u.AddStatusTransition("pin_set", "user")
}

func (u *UserActivationTracking) MarkUserCompleted() error {
	now := time.Now()
	u.UserCompleted = true
	u.UserCompletedAt = &now
	return u.AddStatusTransition("user_registered", "user")
}

func (u *UserActivationTracking) Activate() error {
	now := time.Now()
	u.ActivatedAt = &now

	var method ActivationMethod
	if u.AdminCreatedAt != nil && u.UserCompletedAt != nil {
		if u.AdminCreatedAt.Before(*u.UserCompletedAt) {
			method = ActivationMethodAdminFirst
		} else if u.UserCompletedAt.Before(*u.AdminCreatedAt) {
			method = ActivationMethodUserFirst
		} else {
			method = ActivationMethodSimultaneous
		}
	} else if u.AdminCreatedAt != nil {
		method = ActivationMethodAdminFirst
	} else {
		method = ActivationMethodUserFirst
	}
	u.ActivationMethod = &method

	return u.AddStatusTransition("activated", "system")
}

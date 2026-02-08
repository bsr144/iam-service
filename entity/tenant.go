package entity

import (
	"time"

	"github.com/google/uuid"
)

type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "active"
	TenantStatusInactive  TenantStatus = "inactive"
	TenantStatusSuspended TenantStatus = "suspended"
)

type Tenant struct {
	TenantID     uuid.UUID    `json:"tenant_id" gorm:"column:tenant_id;primaryKey" db:"tenant_id"`
	Name         string       `json:"name" gorm:"column:name" db:"name"`
	Slug         string       `json:"slug" gorm:"column:slug;uniqueIndex" db:"slug"`
	DatabaseName string       `json:"database_name" gorm:"column:database_name;uniqueIndex" db:"database_name"`
	TenantType   *string      `json:"tenant_type,omitempty" gorm:"column:tenant_type" db:"tenant_type"`
	Status       TenantStatus `json:"status" gorm:"column:status;default:active" db:"status"`
	Timestamps
}

func (Tenant) TableName() string {
	return "tenants"
}

type TenantSettings struct {
	TenantSettingID  uuid.UUID `json:"tenant_setting_id" gorm:"column:tenant_setting_id;primaryKey" db:"tenant_setting_id"`
	TenantID         uuid.UUID `json:"tenant_id" gorm:"column:tenant_id;uniqueIndex" db:"tenant_id"`
	SubscriptionTier string    `json:"subscription_tier" gorm:"column:subscription_tier;default:standard" db:"subscription_tier"`
	MaxBranches      int       `json:"max_branches" gorm:"column:max_branches;default:10" db:"max_branches"`
	MaxEmployees     int       `json:"max_employees" gorm:"column:max_employees;default:10000" db:"max_employees"`
	ContactEmail     string    `json:"contact_email,omitempty" gorm:"column:contact_email" db:"contact_email"`
	ContactPhone     string    `json:"contact_phone,omitempty" gorm:"column:contact_phone" db:"contact_phone"`
	ContactAddress   string    `json:"contact_address,omitempty" gorm:"column:contact_address" db:"contact_address"`
	DefaultLanguage  string    `json:"default_language" gorm:"column:default_language;default:en" db:"default_language"`
	Timezone         string    `json:"timezone" gorm:"column:timezone;default:Asia/Jakarta" db:"timezone"`
	CreatedAt        time.Time `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
}

func (TenantSettings) TableName() string {
	return "tenant_settings"
}

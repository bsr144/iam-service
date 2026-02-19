package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ScopeLevel string

const (
	ScopeLevelSystem ScopeLevel = "system"
	ScopeLevelTenant ScopeLevel = "tenant"
	ScopeLevelBranch ScopeLevel = "branch"
	ScopeLevelSelf   ScopeLevel = "self"
)

type Product struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	TenantID    uuid.UUID       `json:"tenant_id" db:"tenant_id"`
	Code        string          `json:"code" db:"code"`
	Name        string          `json:"name" db:"name"`
	Description *string         `json:"description,omitempty" db:"description"`
	Settings    json.RawMessage `json:"settings" db:"settings"`
	Status      string          `json:"status" db:"status"`
	CreatedBy   *uuid.UUID      `json:"created_by,omitempty" db:"created_by"`
	Version     int             `json:"version" db:"version"`
	Timestamps
}

func (Product) TableName() string {
	return "applications"
}

func (p *Product) IsActive() bool {
	return p.Status == "ACTIVE"
}

type Permission struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	ApplicationID uuid.UUID  `json:"application_id" db:"application_id"`
	Code          string     `json:"code" db:"code"`
	Name          string     `json:"name" db:"name"`
	Description   *string    `json:"description,omitempty" db:"description"`
	ResourceType  *string    `json:"resource_type,omitempty" db:"resource_type"`
	Action        *string    `json:"action,omitempty" db:"action"`
	Status        string     `json:"status" db:"status"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty" db:"created_by"`
	Version       int        `json:"version" db:"version"`
	Timestamps
}

func (p *Permission) IsActive() bool {
	return p.Status == "ACTIVE"
}

type Role struct {
	ID            uuid.UUID  `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	ApplicationID uuid.UUID  `json:"application_id" gorm:"column:application_id;not null" db:"application_id"`
	Code          string     `json:"code" db:"code"`
	Name          string     `json:"name" db:"name"`
	Description   *string    `json:"description,omitempty" db:"description"`
	IsSystem      bool       `json:"is_system" db:"is_system"`
	Status        string     `json:"status" db:"status"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty" db:"created_by"`
	Version       int        `json:"version" db:"version"`
	Timestamps
}

func (r *Role) IsActive() bool {
	return r.Status == "ACTIVE"
}

type RolePermission struct {
	ID           uuid.UUID  `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	RoleID       uuid.UUID  `json:"role_id" gorm:"column:role_id;not null" db:"role_id"`
	PermissionID uuid.UUID  `json:"permission_id" gorm:"column:permission_id;not null" db:"permission_id"`
	CreatedBy    *uuid.UUID `json:"created_by,omitempty" gorm:"column:created_by" db:"created_by"`
	CreatedAt    time.Time  `json:"created_at" gorm:"column:created_at" db:"created_at"`
}

type UserRole struct {
	ID         uuid.UUID  `json:"id" gorm:"column:id;primaryKey;type:uuid;default:uuidv7()" db:"id"`
	UserID     uuid.UUID  `json:"user_id" gorm:"column:user_id;type:uuid;not null" db:"user_id"`
	RoleID     uuid.UUID  `json:"role_id" gorm:"column:role_id;type:uuid;not null" db:"role_id"`
	ProductID  *uuid.UUID `json:"product_id,omitempty" gorm:"column:product_id;type:uuid" db:"product_id"`
	BranchID   *uuid.UUID `json:"branch_id,omitempty" gorm:"column:branch_id;type:uuid" db:"branch_id"`
	AssignedAt time.Time  `json:"assigned_at" gorm:"column:assigned_at;not null" db:"assigned_at"`
	AssignedBy *uuid.UUID `json:"assigned_by,omitempty" gorm:"column:assigned_by" db:"assigned_by"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty" gorm:"column:expires_at" db:"expires_at"`
	Status     string     `json:"status" gorm:"column:status" db:"status"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"column:updated_at" db:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at" db:"deleted_at"`
}

func (UserRole) TableName() string {
	return "user_role_assignments"
}

func (ur *UserRole) IsActive() bool {
	if ur.DeletedAt != nil {
		return false
	}
	if ur.Status != "ACTIVE" {
		return false
	}
	if ur.ExpiresAt != nil && time.Now().After(*ur.ExpiresAt) {
		return false
	}
	return true
}

func (ur *UserRole) IsTenantWide() bool {
	return ur.BranchID == nil
}

func (ur *UserRole) IsProductSpecific() bool {
	return ur.ProductID != nil
}

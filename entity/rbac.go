package entity

import (
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
	ProductID     uuid.UUID  `json:"product_id" db:"product_id"`
	TenantID      uuid.UUID  `json:"tenant_id" db:"tenant_id"`
	Code          string     `json:"code" db:"code"`
	Name          string     `json:"name" db:"name"`
	Description   *string    `json:"description,omitempty" db:"description"`
	ProductType   string     `json:"product_type" db:"product_type"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	LicensedUntil *time.Time `json:"licensed_until,omitempty" db:"licensed_until"`
	Timestamps
}

func (p *Product) IsLicensed() bool {
	if !p.IsActive {
		return false
	}
	if p.LicensedUntil == nil {
		return true
	}
	return time.Now().Before(*p.LicensedUntil)
}

type Permission struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Code        string     `json:"code" db:"code"`
	Name        string     `json:"name" db:"name"`
	Description *string    `json:"description,omitempty" db:"description"`
	Module      string     `json:"module" db:"module"`
	Resource    string     `json:"resource" db:"resource"`
	Action      string     `json:"action" db:"action"`
	ScopeLevel  ScopeLevel `json:"scope_level" db:"scope_level"`
	IsSystem    bool       `json:"is_system" db:"is_system"`
	Timestamps
}

type Role struct {
	RoleID       uuid.UUID  `json:"role_id" db:"role_id"`
	TenantID     *uuid.UUID `json:"tenant_id,omitempty" db:"tenant_id"`
	ProductID    *uuid.UUID `json:"product_id,omitempty" db:"product_id"`
	Code         string     `json:"code" db:"code"`
	Name         string     `json:"name" db:"name"`
	Description  *string    `json:"description,omitempty" db:"description"`
	ParentRoleID *uuid.UUID `json:"parent_role_id,omitempty" db:"parent_role_id"`
	ScopeLevel   ScopeLevel `json:"scope_level" db:"scope_level"`
	IsSystem     bool       `json:"is_system" db:"is_system"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	Timestamps
}

func (r *Role) IsSystemRole() bool {
	return r.TenantID == nil
}

func (r *Role) IsTenantWide() bool {
	return r.TenantID != nil && r.ProductID == nil
}

func (r *Role) IsProductSpecific() bool {
	return r.ProductID != nil
}

type RolePermission struct {
	RolePermissionID uuid.UUID `json:"role_permission_id" gorm:"column:role_permission_id;primaryKey" db:"role_permission_id"`
	RoleID           uuid.UUID `json:"role_id" gorm:"column:role_id;not null" db:"role_id"`
	PermissionID     uuid.UUID `json:"permission_id" gorm:"column:permission_id;not null" db:"permission_id"`
	CreatedAt        time.Time `json:"created_at" gorm:"column:created_at" db:"created_at"`
}

type UserRole struct {
	UserRoleID    uuid.UUID  `json:"user_role_id" db:"user_role_id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	RoleID        uuid.UUID  `json:"role_id" db:"role_id"`
	ProductID     *uuid.UUID `json:"product_id,omitempty" db:"product_id"`
	BranchID      *uuid.UUID `json:"branch_id,omitempty" db:"branch_id"`
	EffectiveFrom time.Time  `json:"effective_from" db:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty" db:"effective_to"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

func (ur *UserRole) IsActive() bool {
	now := time.Now()
	if ur.DeletedAt != nil {
		return false
	}
	if now.Before(ur.EffectiveFrom) {
		return false
	}
	if ur.EffectiveTo != nil && now.After(*ur.EffectiveTo) {
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

package userdto

import "github.com/google/uuid"

type CreateRequest struct {
	TenantID  uuid.UUID  `json:"tenant_id" validate:"required"`
	RoleCode  string     `json:"role_code" validate:"required"`
	Email     string     `json:"email" validate:"required,email,max=255"`
	Password  string     `json:"password" validate:"required,min=8,max=128"`
	FirstName string     `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string     `json:"last_name" validate:"required,min=2,max=100"`
	BranchID  *uuid.UUID `json:"branch_id,omitempty" validate:"omitempty,uuid"`
}

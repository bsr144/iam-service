package response

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID          `json:"id"`
	Email     string             `json:"email"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	FullName  string             `json:"full_name"`
	Phone     *string            `json:"phone,omitempty"`
	IsActive  bool               `json:"is_active"`
	TenantID  *uuid.UUID         `json:"tenant_id,omitempty"`
	BranchID  *uuid.UUID         `json:"branch_id,omitempty"`
	Roles     []UserRoleResponse `json:"roles,omitempty"`
}

type UserRoleResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type UserListItemResponse struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	FullName  string     `json:"full_name"`
	Phone     *string    `json:"phone,omitempty"`
	IsActive  bool       `json:"is_active"`
	TenantID  *uuid.UUID `json:"tenant_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type CreateUserResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	RoleCode string    `json:"role_code"`
}

type UpdateUserResponse struct {
	UserResponse
	Message string `json:"message"`
}

type ApproveUserResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message"`
}

type RejectUserResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message"`
}

type UnlockUserResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message"`
}

type ResetUserPINResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message"`
}

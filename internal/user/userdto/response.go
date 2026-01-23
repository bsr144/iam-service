package userdto

import "github.com/google/uuid"

type CreateResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	RoleCode string    `json:"role_code"`
	TenantID uuid.UUID `json:"tenant_id"`
}

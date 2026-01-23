package entity

import (
	"time"

	"github.com/google/uuid"
)

type Branch struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	TenantID      uuid.UUID  `json:"tenant_id" db:"tenant_id"`
	Name          string     `json:"name" db:"name"`
	Code          string     `json:"code" db:"code"`
	IsHeadquarters bool      `json:"is_headquarters" db:"is_headquarters"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	Timestamps
}

type BranchContact struct {
	ID         uuid.UUID `json:"id" db:"id"`
	BranchID   uuid.UUID `json:"branch_id" db:"branch_id"`
	Address    string    `json:"address,omitempty" db:"address"`
	City       string    `json:"city,omitempty" db:"city"`
	Province   string    `json:"province,omitempty" db:"province"`
	PostalCode string    `json:"postal_code,omitempty" db:"postal_code"`
	Phone      string    `json:"phone,omitempty" db:"phone"`
	Email      string    `json:"email,omitempty" db:"email"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

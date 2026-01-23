package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID      uuid.UUID  `json:"user_id"`
	Email       string     `json:"email"`
	TenantID    *uuid.UUID `json:"tenant_id,omitempty"`
	ProductID   *uuid.UUID `json:"product_id,omitempty"`
	Roles       []string   `json:"roles"`
	Permissions []string   `json:"permissions,omitempty"`
	BranchID    *uuid.UUID `json:"branch_id,omitempty"`
	SessionID   uuid.UUID  `json:"session_id"`
	jwt.RegisteredClaims
}

func (c *JWTClaims) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return c.ExpiresAt.Before(time.Now())
}

func (c *JWTClaims) IsPlatformAdmin() bool {
	if c.TenantID != nil {
		return false
	}
	for _, role := range c.Roles {
		if role == "PLATFORM_ADMIN" {
			return true
		}
	}
	return false
}

func (c *JWTClaims) HasRole(roleCode string) bool {
	for _, role := range c.Roles {
		if role == roleCode {
			return true
		}
	}
	return false
}

func (c *JWTClaims) IsTenantUser() bool {
	return c.TenantID != nil
}

func (c *JWTClaims) GetTenantID() uuid.UUID {
	if c.TenantID == nil {
		return uuid.Nil
	}
	return *c.TenantID
}

func (c *JWTClaims) GetBranchID() uuid.UUID {
	if c.BranchID == nil {
		return uuid.Nil
	}
	return *c.BranchID
}

func (c *JWTClaims) GetProductID() uuid.UUID {
	if c.ProductID == nil {
		return uuid.Nil
	}
	return *c.ProductID
}

func (c *JWTClaims) HasProductContext() bool {
	return c.ProductID != nil
}

func (c *JWTClaims) HasPermission(permissionCode string) bool {
	for _, perm := range c.Permissions {
		if perm == permissionCode {
			return true
		}
	}
	return false
}

func (c *JWTClaims) HasAudience(audience string) bool {
	for _, aud := range c.Audience {
		if aud == audience {
			return true
		}
	}
	return false
}

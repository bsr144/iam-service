package middleware

import (
	"iam-service/pkg/errors"
	jwtpkg "iam-service/pkg/jwt"
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	UserClaimsKey = "user_claims"
)

func GetUserClaims(c *fiber.Ctx) (*jwtpkg.JWTClaims, error) {
	claims := c.Locals(UserClaimsKey)
	if claims == nil {
		return nil, errors.ErrUnauthorized("user claims not found in context")
	}

	jwtClaims, ok := claims.(*jwtpkg.JWTClaims)
	if !ok {
		return nil, errors.ErrInternal("invalid claims type in context")
	}

	return jwtClaims, nil
}

func GetUserID(c *fiber.Ctx) (uuid.UUID, error) {
	claims, err := GetUserClaims(c)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.UserID, nil
}

func GetTenantID(c *fiber.Ctx) (*uuid.UUID, error) {
	claims, err := GetUserClaims(c)
	if err != nil {
		return nil, err
	}
	return claims.TenantID, nil
}

func GetBranchID(c *fiber.Ctx) (*uuid.UUID, error) {
	claims, err := GetUserClaims(c)
	if err != nil {
		return nil, err
	}
	return claims.BranchID, nil
}

func GetSessionID(c *fiber.Ctx) (uuid.UUID, error) {
	claims, err := GetUserClaims(c)
	if err != nil {
		return uuid.Nil, err
	}
	return claims.SessionID, nil
}

func GetUserRoles(c *fiber.Ctx) ([]string, error) {
	claims, err := GetUserClaims(c)
	if err != nil {
		return nil, err
	}
	return claims.Roles, nil
}

func GetClientIP(c *fiber.Ctx) net.IP {

	forwarded := c.Get("X-Forwarded-For")
	if forwarded != "" {

		return net.ParseIP(forwarded)
	}

	return net.ParseIP(c.IP())
}

func GetUserAgent(c *fiber.Ctx) string {
	return c.Get("User-Agent")
}

func GetRequestID(c *fiber.Ctx) string {
	if id := c.GetRespHeader("X-Request-ID"); id != "" {
		return id
	}
	return c.Get("X-Request-ID")
}

func IsPlatformAdmin(c *fiber.Ctx) (bool, error) {
	claims, err := GetUserClaims(c)
	if err != nil {
		return false, err
	}
	return claims.IsPlatformAdmin(), nil
}

func GetTenantIDFromHeader(c *fiber.Ctx) (uuid.UUID, error) {
	tenantIDStr := c.Get("X-Tenant-ID")
	if tenantIDStr == "" {
		return uuid.Nil, errors.ErrBadRequest("X-Tenant-ID header is required")
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return uuid.Nil, errors.ErrBadRequest("Invalid X-Tenant-ID format")
	}

	return tenantID, nil
}

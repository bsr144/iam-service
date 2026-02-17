package middleware

import (
	"iam-service/pkg/errors"
	jwtpkg "iam-service/pkg/jwt"
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	UserClaimsKey         = "user_claims"
	MultiTenantClaimsKey  = "multi_tenant_claims"
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

func GetMultiTenantClaims(c *fiber.Ctx) (*jwtpkg.MultiTenantClaims, error) {
	claims := c.Locals(MultiTenantClaimsKey)
	if claims == nil {
		return nil, errors.ErrUnauthorized("multi-tenant claims not found in context")
	}

	multiClaims, ok := claims.(*jwtpkg.MultiTenantClaims)
	if !ok {
		return nil, errors.ErrInternal("invalid multi-tenant claims type in context")
	}

	return multiClaims, nil
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

func ExtractTenantContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		multiClaims, err := GetMultiTenantClaims(c)
		if err != nil {
			appErr := errors.ErrUnauthorized("authentication required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		tenantID, err := GetTenantIDFromHeader(c)
		if err != nil {
			appErr := errors.ErrBadRequest("X-Tenant-ID header is required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		if !multiClaims.HasTenant(tenantID) {
			appErr := errors.ErrForbidden("access denied to this tenant")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		c.Locals("tenant_id", tenantID)

		return c.Next()
	}
}

func GetTenantIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	tenantID := c.Locals("tenant_id")
	if tenantID == nil {
		return uuid.Nil, errors.ErrForbidden("tenant context not found")
	}

	tid, ok := tenantID.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.ErrInternal("invalid tenant ID type in context")
	}

	return tid, nil
}

func ExtractProductContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		productIDStr := c.Params("productId")
		if productIDStr == "" {
			appErr := errors.ErrBadRequest("product ID is required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		productID, err := uuid.Parse(productIDStr)
		if err != nil {
			appErr := errors.ErrBadRequest("invalid product ID format")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		c.Locals("product_id", productID)

		return c.Next()
	}
}

func GetProductIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	productID := c.Locals("product_id")
	if productID == nil {
		return uuid.Nil, errors.ErrBadRequest("product context not found")
	}

	pid, ok := productID.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.ErrInternal("invalid product ID type in context")
	}

	return pid, nil
}

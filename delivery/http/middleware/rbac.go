package middleware

import (
	"iam-service/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

func RequirePlatformAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := GetUserClaims(c)
		if err != nil {

			appErr := errors.ErrUnauthorized("authentication required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		if !claims.IsPlatformAdmin() {
			appErr := errors.ErrPlatformAdminRequired()
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		return c.Next()
	}
}
func RequireRole(roleCode string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := GetUserClaims(c)
		if err != nil {
			appErr := errors.ErrUnauthorized("authentication required")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		if !claims.HasRole(roleCode) {
			appErr := errors.ErrAccessForbidden("insufficient permissions")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		return c.Next()
	}
}

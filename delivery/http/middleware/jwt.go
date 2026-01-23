package middleware

import (
	"iam-service/config"
	"iam-service/pkg/errors"
	jwtpkg "iam-service/pkg/jwt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTAuth(cfg *config.Config) fiber.Handler {
	tokenConfig := &jwtpkg.TokenConfig{
		AccessSecret:  cfg.JWT.AccessSecret,
		RefreshSecret: cfg.JWT.RefreshSecret,
		AccessExpiry:  cfg.JWT.AccessExpiry,
		RefreshExpiry: cfg.JWT.RefreshExpiry,
		Issuer:        cfg.JWT.Issuer,
	}

	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			appErr := errors.ErrUnauthorized("missing authorization header")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			appErr := errors.ErrUnauthorized("invalid authorization header format")
			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		tokenString := parts[1]

		claims, err := jwtpkg.ParseAccessToken(tokenString, tokenConfig)
		if err != nil {
			var appErr *errors.AppError
			switch err {
			case jwtpkg.ErrTokenExpired:
				appErr = errors.ErrTokenExpired()
			case jwtpkg.ErrTokenInvalid, jwtpkg.ErrTokenMalformed, jwtpkg.ErrTokenSignature:
				appErr = errors.ErrTokenInvalid()
			default:
				appErr = errors.ErrTokenInvalid()
			}

			return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
				"success": false,
				"error":   appErr.Message,
				"code":    appErr.Code,
			})
		}

		c.Locals(UserClaimsKey, claims)

		return c.Next()
	}
}

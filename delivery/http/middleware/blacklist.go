package middleware

import (
	"iam-service/iam/auth/contract"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// checkBlacklist checks if a token JTI or user is blacklisted.
// Returns a fiber error response if blocked, nil if allowed.
// Fail-open: Redis errors are silently ignored (token remains cryptographically valid).
func checkBlacklist(c *fiber.Ctx, store contract.TokenBlacklistStore, jti string, userID uuid.UUID, claims jwt.RegisteredClaims) error {
	if jti == "" {
		return nil
	}

	// Check per-token blacklist
	blacklisted, err := store.IsTokenBlacklisted(c.UserContext(), jti)
	if err == nil && blacklisted {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "token has been revoked",
			"code":    "TOKEN_REVOKED",
		})
	}

	// Check user-level blacklist (logout-all)
	blacklistTS, err := store.GetUserBlacklistTimestamp(c.UserContext(), userID)
	if err == nil && blacklistTS != nil && claims.IssuedAt != nil {
		if claims.IssuedAt.Time.Before(*blacklistTS) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "token has been revoked",
				"code":    "TOKEN_REVOKED",
			})
		}
	}

	return nil
}

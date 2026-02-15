package contract

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// TokenBlacklistStore defines operations for blacklisting access tokens and users
// Implemented by: impl/redis/token_blacklist_store.go
type TokenBlacklistStore interface {
	// BlacklistToken adds an access token's JTI to the blacklist
	// ttl should match the token's remaining lifetime (exp - now)
	BlacklistToken(ctx context.Context, jti string, ttl time.Duration) error

	// IsTokenBlacklisted checks if a token JTI is blacklisted
	IsTokenBlacklisted(ctx context.Context, jti string) (bool, error)

	// BlacklistUser adds a user to the user-level blacklist with a timestamp
	// All tokens issued before this timestamp are invalidated
	// ttl should be the maximum token lifetime (e.g., 30 days for refresh tokens)
	BlacklistUser(ctx context.Context, userID uuid.UUID, timestamp time.Time, ttl time.Duration) error

	// GetUserBlacklistTimestamp retrieves the user blacklist timestamp
	// Returns nil if user is not blacklisted
	GetUserBlacklistTimestamp(ctx context.Context, userID uuid.UUID) (*time.Time, error)
}

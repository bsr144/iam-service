package redis

import (
	"fmt"
)

// Error factory functions for wrapping errors with context.
// Sentinel errors are defined in pkg/errors/sentinel.go.
var (
	ErrFailedToMarshalValue = func(err error) error {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	ErrFailedToCheckRateLimit = func(err error) error {
		return fmt.Errorf("failed to check rate limit: %w", err)
	}
)

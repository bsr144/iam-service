package redis

import (
	"errors"
	"fmt"
)

var (
	ErrFailedToMarshalValue = func(err error) error {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	ErrFailedToCheckRateLimit = func(err error) error {
		return fmt.Errorf("failed to check rate limit: %w", err)
	}

	ErrLockNotAcquired        = errors.New("lock not acquired")
	ErrLockNotHeld            = errors.New("lock not held")
	ErrSemaphoreFull          = errors.New("semaphore is full")
	ErrSemaphoreTokenNotFound = errors.New("semaphore token not found")
)

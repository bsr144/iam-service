package errors

import "errors"

// Sentinel errors for infrastructure operations where callers need to
// distinguish between "expected state" (e.g., lock busy) and "actual error".
// These are used internally by impl packages for operational states.
var (
	// Distributed lock sentinel errors
	SentinelLockNotAcquired = errors.New("lock not acquired")
	SentinelLockNotHeld     = errors.New("lock not held")

	// Semaphore sentinel errors
	SentinelSemaphoreFull          = errors.New("semaphore is full")
	SentinelSemaphoreTokenNotFound = errors.New("semaphore token not found")

	// Queue sentinel errors
	SentinelQueueEmpty  = errors.New("queue is empty")
	SentinelJobNotFound = errors.New("job not found")
)

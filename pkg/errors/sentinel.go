package errors

import "errors"

var (
	// Database sentinel errors
	SentinelNotFound   = errors.New("record not found")
	SentinelDuplicate  = errors.New("duplicate entry")
	SentinelForeignKey = errors.New("foreign key violation")
	SentinelConflict   = errors.New("optimistic lock conflict")

	// Cache sentinel errors
	SentinelCacheMiss    = errors.New("cache miss")
	SentinelCacheTimeout = errors.New("cache timeout")

	// Distributed lock sentinel errors
	SentinelLockNotAcquired = errors.New("lock not acquired")
	SentinelLockNotHeld     = errors.New("lock not held")

	// Semaphore sentinel errors
	SentinelSemaphoreFull          = errors.New("semaphore is full")
	SentinelSemaphoreTokenNotFound = errors.New("semaphore token not found")

	// Object storage sentinel errors
	SentinelObjectNotFound = errors.New("object not found")
	SentinelBucketNotFound = errors.New("bucket not found")
	SentinelQuotaExceeded  = errors.New("storage quota exceeded")

	// Secret management sentinel errors
	SentinelSecretNotFound = errors.New("secret not found")

	// Queue sentinel errors
	SentinelQueueEmpty  = errors.New("queue is empty")
	SentinelJobNotFound = errors.New("job not found")
)

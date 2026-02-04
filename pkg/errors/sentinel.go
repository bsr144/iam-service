package errors

import "errors"

var (
	SentinelNotFound   = errors.New("record not found")
	SentinelDuplicate  = errors.New("duplicate entry")
	SentinelForeignKey = errors.New("foreign key violation")
	SentinelConflict   = errors.New("optimistic lock conflict")

	SentinelCacheMiss    = errors.New("cache miss")
	SentinelCacheTimeout = errors.New("cache timeout")

	SentinelObjectNotFound = errors.New("object not found")
	SentinelBucketNotFound = errors.New("bucket not found")
	SentinelQuotaExceeded  = errors.New("storage quota exceeded")
)

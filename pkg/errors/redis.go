package errors

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

func TranslateRedis(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, redis.Nil) {
		return SentinelCacheMiss
	}

	return err
}

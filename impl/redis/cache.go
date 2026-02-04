package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pkgerrors "iam-service/pkg/errors"
)

func (c *Redis) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return ErrFailedToMarshalValue(err)
	}
	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *Redis) SetString(ctx context.Context, key, value string, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *Redis) Get(ctx context.Context, key string, target any) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if translated := pkgerrors.TranslateRedis(err); translated != err {
			return translated
		}
		return fmt.Errorf("failed to get value: %w", err)
	}
	return json.Unmarshal(data, target)
}

func (c *Redis) GetString(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if translated := pkgerrors.TranslateRedis(err); translated != err {
			return "", translated
		}
		return "", fmt.Errorf("failed to get value: %w", err)
	}
	return val, nil
}

func (c *Redis) Delete(ctx context.Context, keys ...string) error {
	prefixedKeys := make([]string, len(keys))
	copy(prefixedKeys, keys)
	return c.client.Del(ctx, prefixedKeys...).Err()
}

func (c *Redis) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return count > 0, nil
}

func (c *Redis) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

func (c *Redis) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL: %w", err)
	}
	return ttl, nil
}

func (c *Redis) Increment(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

func (c *Redis) IncrementBy(ctx context.Context, key string, delta int64) (int64, error) {
	return c.client.IncrBy(ctx, key, delta).Result()
}

func (c *Redis) Decrement(ctx context.Context, key string) (int64, error) {
	return c.client.Decr(ctx, key).Result()
}

func (c *Redis) DecrementBy(ctx context.Context, key string, delta int64) (int64, error) {
	return c.client.DecrBy(ctx, key, delta).Result()
}

func (c *Redis) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}
	return c.client.SetNX(ctx, key, data, expiration).Result()
}

func (c *Redis) GetOrSet(ctx context.Context, key string, target any, expiration time.Duration, fn func() (any, error)) error {
	err := c.Get(ctx, key, target)
	if err == nil {
		return nil
	}
	if err != pkgerrors.SentinelCacheMiss {
		return err
	}

	value, err := fn()
	if err != nil {
		return err
	}

	if err := c.Set(ctx, key, value, expiration); err != nil {
		return err
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func (c *Redis) DeleteByPattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan keys: %w", err)
	}
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}

func (c *Redis) HSet(ctx context.Context, key, field string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return ErrFailedToMarshalValue(err)
	}
	return c.client.HSet(ctx, key, field, data).Err()
}

func (c *Redis) HGet(ctx context.Context, key, field string, target any) error {
	data, err := c.client.HGet(ctx, key, field).Bytes()
	if err != nil {
		if translated := pkgerrors.TranslateRedis(err); translated != err {
			return translated
		}
		return fmt.Errorf("failed to get hash field: %w", err)
	}
	return json.Unmarshal(data, target)
}

func (c *Redis) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	result, err := c.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get hash: %w", err)
	}
	return result, nil
}

func (c *Redis) HDel(ctx context.Context, key string, fields ...string) error {
	return c.client.HDel(ctx, key, fields...).Err()
}

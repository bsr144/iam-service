package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Lock struct {
	redis    *Redis
	name     string
	token    string
	expiry   time.Duration
	acquired bool
}

func generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (r *Redis) AcquireLock(ctx context.Context, name string, expiry time.Duration) (*Lock, error) {
	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate lock token: %w", err)
	}

	lock := &Lock{
		redis:  r,
		name:   name,
		token:  token,
		expiry: expiry,
	}

	lockKey := fmt.Sprintf(LockPrefix, name)
	acquired, err := r.client.SetNX(ctx, lockKey, token, expiry).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}

	if !acquired {
		return nil, ErrLockNotAcquired
	}

	lock.acquired = true
	return lock, nil
}

func (r *Redis) AcquireLockWithRetry(ctx context.Context, name string, expiry time.Duration, maxRetries int, retryDelay time.Duration) (*Lock, error) {
	for i := 0; i < maxRetries; i++ {
		lock, err := r.AcquireLock(ctx, name, expiry)
		if err == nil {
			return lock, nil
		}
		if !errors.Is(err, ErrLockNotAcquired) {
			return nil, err
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(retryDelay):
			continue
		}
	}
	return nil, ErrLockNotAcquired
}

func (r *Redis) AcquireLockWithWait(ctx context.Context, name string, expiry time.Duration, pollInterval time.Duration) (*Lock, error) {
	for {
		lock, err := r.AcquireLock(ctx, name, expiry)
		if err == nil {
			return lock, nil
		}
		if !errors.Is(err, ErrLockNotAcquired) {
			return nil, err
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(pollInterval):
			continue
		}
	}
}

func (lock *Lock) Release(ctx context.Context) error {
	if !lock.acquired {
		return ErrLockNotHeld
	}

	script := redis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`)

	lockKey := fmt.Sprintf(LockPrefix, lock.name)
	result, err := script.Run(ctx, lock.redis.client, []string{lockKey}, lock.token).Int64()
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	if result == 0 {
		return ErrLockNotHeld
	}

	lock.acquired = false
	return nil
}

func (lock *Lock) Extend(ctx context.Context, expiry time.Duration) error {
	if !lock.acquired {
		return ErrLockNotHeld
	}

	script := redis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("pexpire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`)

	lockKey := fmt.Sprintf(LockPrefix, lock.name)
	result, err := script.Run(ctx, lock.redis.client, []string{lockKey}, lock.token, int64(expiry/time.Millisecond)).Int64()
	if err != nil {
		return fmt.Errorf("failed to extend lock: %w", err)
	}

	if result == 0 {
		lock.acquired = false
		return ErrLockNotHeld
	}

	lock.expiry = expiry
	return nil
}

func (lock *Lock) IsHeld(ctx context.Context) (bool, error) {
	lockKey := fmt.Sprintf(LockPrefix, lock.name)
	val, err := lock.redis.client.Get(ctx, lockKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			lock.acquired = false
			return false, nil
		}
		return false, fmt.Errorf("failed to check lock: %w", err)
	}
	held := val == lock.token
	if !held {
		lock.acquired = false
	}
	return held, nil
}

func (lock *Lock) TTL(ctx context.Context) (time.Duration, error) {
	lockKey := fmt.Sprintf(LockPrefix, lock.name)
	ttl, err := lock.redis.client.TTL(ctx, lockKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get lock TTL: %w", err)
	}
	return ttl, nil
}

func (r *Redis) WithLock(ctx context.Context, name string, expiry time.Duration, fn func(ctx context.Context) error) error {
	lock, err := r.AcquireLock(ctx, name, expiry)
	if err != nil {
		return err
	}
	defer lock.Release(ctx)

	return fn(ctx)
}

func (r *Redis) WithLockRetry(ctx context.Context, name string, expiry time.Duration, maxRetries int, retryDelay time.Duration, fn func(ctx context.Context) error) error {
	lock, err := r.AcquireLockWithRetry(ctx, name, expiry, maxRetries, retryDelay)
	if err != nil {
		return err
	}
	defer lock.Release(ctx)

	return fn(ctx)
}

func (r *Redis) IsLocked(ctx context.Context, name string) (bool, error) {
	lockKey := fmt.Sprintf(LockPrefix, name)
	exists, err := r.client.Exists(ctx, lockKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check lock: %w", err)
	}
	return exists > 0, nil
}

func (r *Redis) ForceReleaseLock(ctx context.Context, name string) error {
	lockKey := fmt.Sprintf(LockPrefix, name)
	return r.client.Del(ctx, lockKey).Err()
}

func (r *Redis) NewSemaphore(name string, maxCount int64) *Semaphore {
	return &Semaphore{
		redis:    r,
		name:     name,
		maxCount: maxCount,
	}
}

func (s *Semaphore) Acquire(ctx context.Context, expiry time.Duration) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	script := redis.NewScript(`
		local current = redis.call("ZCARD", KEYS[1])
		if current < tonumber(ARGV[1]) then
			redis.call("ZADD", KEYS[1], ARGV[2], ARGV[3])
			return 1
		else
			return 0
		end
	`)

	semKey := fmt.Sprintf(SemaphorePrefix, s.name)
	expireAt := float64(time.Now().Add(expiry).Unix())
	result, err := script.Run(ctx, s.redis.client, []string{semKey}, s.maxCount, expireAt, token).Int64()
	if err != nil {
		return "", fmt.Errorf("failed to acquire semaphore: %w", err)
	}

	if result == 0 {
		return "", ErrSemaphoreFull
	}

	return token, nil
}

func (s *Semaphore) Release(ctx context.Context, token string) error {
	semKey := fmt.Sprintf(SemaphorePrefix, s.name)
	removed, err := s.redis.client.ZRem(ctx, semKey, token).Result()
	if err != nil {
		return fmt.Errorf("failed to release semaphore: %w", err)
	}
	if removed == 0 {
		return ErrSemaphoreTokenNotFound
	}
	return nil
}

func (s *Semaphore) Cleanup(ctx context.Context) error {
	semKey := fmt.Sprintf(SemaphorePrefix, s.name)
	now := float64(time.Now().Unix())
	return s.redis.client.ZRemRangeByScore(ctx, semKey, "-inf", fmt.Sprintf("%f", now)).Err()
}

func (s *Semaphore) Count(ctx context.Context) (int64, error) {
	semKey := fmt.Sprintf(SemaphorePrefix, s.name)
	return s.redis.client.ZCard(ctx, semKey).Result()
}

func (s *Semaphore) Available(ctx context.Context) (int64, error) {
	count, err := s.Count(ctx)
	if err != nil {
		return 0, err
	}
	return s.maxCount - count, nil
}

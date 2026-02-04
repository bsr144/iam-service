package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"iam-service/entity"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

const (
	registrationSessionPrefix = "reg:%s:%s"
	registrationEmailPrefix   = "reg_email:%s:%s"
	registrationRatePrefix    = "reg_rate:%s:%s"
)

func (r *Redis) registrationSessionKey(tenantID, sessionID uuid.UUID) string {
	return fmt.Sprintf(registrationSessionPrefix, tenantID.String(), sessionID.String())
}

func (r *Redis) registrationEmailLockKey(tenantID uuid.UUID, email string) string {
	return fmt.Sprintf(registrationEmailPrefix, tenantID.String(), strings.ToLower(email))
}

func (r *Redis) registrationRateLimitKey(tenantID uuid.UUID, email string) string {
	return fmt.Sprintf(registrationRatePrefix, tenantID.String(), strings.ToLower(email))
}

func (r *Redis) CreateRegistrationSession(ctx context.Context, session *entity.RegistrationSession, ttl time.Duration) error {
	key := r.registrationSessionKey(session.TenantID, session.ID)

	data, err := json.Marshal(session)
	if err != nil {
		return errors.ErrInternal("failed to marshal registration session").WithError(err)
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return errors.ErrInternal("failed to store registration session").WithError(err)
	}

	return nil
}

func (r *Redis) GetRegistrationSession(ctx context.Context, tenantID, sessionID uuid.UUID) (*entity.RegistrationSession, error) {
	key := r.registrationSessionKey(tenantID, sessionID)

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, errors.ErrNotFound("registration session not found or expired")
		}
		return nil, errors.ErrInternal("failed to get registration session").WithError(err)
	}

	var session entity.RegistrationSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, errors.ErrInternal("failed to unmarshal registration session").WithError(err)
	}

	return &session, nil
}

func (r *Redis) UpdateRegistrationSession(ctx context.Context, session *entity.RegistrationSession, ttl time.Duration) error {
	key := r.registrationSessionKey(session.TenantID, session.ID)

	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return errors.ErrInternal("failed to check session existence").WithError(err)
	}
	if exists == 0 {
		return errors.ErrNotFound("registration session not found or expired")
	}

	data, err := json.Marshal(session)
	if err != nil {
		return errors.ErrInternal("failed to marshal registration session").WithError(err)
	}

	if ttl == 0 {
		remainingTTL, err := r.client.TTL(ctx, key).Result()
		if err != nil {
			return errors.ErrInternal("failed to get session TTL").WithError(err)
		}
		if remainingTTL > 0 {
			ttl = remainingTTL
		}
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return errors.ErrInternal("failed to update registration session").WithError(err)
	}

	return nil
}

func (r *Redis) DeleteRegistrationSession(ctx context.Context, tenantID, sessionID uuid.UUID) error {
	key := r.registrationSessionKey(tenantID, sessionID)
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) IncrementRegistrationAttempts(ctx context.Context, tenantID, sessionID uuid.UUID) (int, error) {
	session, err := r.GetRegistrationSession(ctx, tenantID, sessionID)
	if err != nil {
		return 0, err
	}

	session.Attempts++

	if session.Attempts >= session.MaxAttempts {
		session.Status = entity.RegistrationSessionStatusFailed
	}

	if err := r.UpdateRegistrationSession(ctx, session, 0); err != nil {
		return 0, err
	}

	return session.Attempts, nil
}

func (r *Redis) UpdateRegistrationOTP(ctx context.Context, tenantID, sessionID uuid.UUID, otpHash string, expiresAt time.Time) error {
	session, err := r.GetRegistrationSession(ctx, tenantID, sessionID)
	if err != nil {
		return err
	}

	now := time.Now()
	session.OTPHash = otpHash
	session.OTPCreatedAt = now
	session.OTPExpiresAt = expiresAt
	session.ResendCount++
	session.LastResentAt = &now

	return r.UpdateRegistrationSession(ctx, session, 0)
}

func (r *Redis) MarkRegistrationVerified(ctx context.Context, tenantID, sessionID uuid.UUID, tokenHash string) error {
	session, err := r.GetRegistrationSession(ctx, tenantID, sessionID)
	if err != nil {
		return err
	}

	now := time.Now()
	session.Status = entity.RegistrationSessionStatusVerified
	session.VerifiedAt = &now
	session.RegistrationTokenHash = &tokenHash

	return r.UpdateRegistrationSession(ctx, session, 0)
}

func (r *Redis) LockRegistrationEmail(ctx context.Context, tenantID uuid.UUID, email string, ttl time.Duration) (bool, error) {
	key := r.registrationEmailLockKey(tenantID, email)
	return r.client.SetNX(ctx, key, "1", ttl).Result()
}

func (r *Redis) UnlockRegistrationEmail(ctx context.Context, tenantID uuid.UUID, email string) error {
	key := r.registrationEmailLockKey(tenantID, email)
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) IsRegistrationEmailLocked(ctx context.Context, tenantID uuid.UUID, email string) (bool, error) {
	key := r.registrationEmailLockKey(tenantID, email)
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, errors.ErrInternal("failed to check email lock").WithError(err)
	}
	return exists > 0, nil
}

func (r *Redis) IncrementRegistrationRateLimit(ctx context.Context, tenantID uuid.UUID, email string, ttl time.Duration) (int64, error) {
	key := r.registrationRateLimitKey(tenantID, email)

	pipe := r.client.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, ttl)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, errors.ErrInternal("failed to increment rate limit").WithError(err)
	}

	return incr.Val(), nil
}

func (r *Redis) GetRegistrationRateLimitCount(ctx context.Context, tenantID uuid.UUID, email string) (int64, error) {
	key := r.registrationRateLimitKey(tenantID, email)

	count, err := r.client.Get(ctx, key).Int64()
	if err != nil {
		if err == goredis.Nil {
			return 0, nil
		}
		return 0, errors.ErrInternal("failed to get rate limit count").WithError(err)
	}

	return count, nil
}

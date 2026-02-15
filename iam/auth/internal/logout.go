package internal

import (
	"context"
	"time"

	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"
)

func (uc *usecase) Logout(ctx context.Context, req *authdto.LogoutRequest) error {
	if req.RefreshToken == "" {
		return errors.ErrBadRequest("refresh token is required")
	}

	tokenHash := hashToken(req.RefreshToken)

	refreshToken, err := uc.RefreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil // idempotent: not-found is OK
		}
		return errors.ErrInternal("failed to verify token").WithError(err)
	}

	// BOLA check: verify ownership
	if refreshToken.UserID != req.UserID {
		return nil // idempotent — don't reveal token exists for other user
	}

	if refreshToken.RevokedAt != nil {
		return nil // idempotent: already revoked is OK
	}

	if refreshToken.IsExpired() {
		return nil // idempotent: expired is OK
	}

	// Find linked session
	session, _ := uc.UserSessionRepo.GetByRefreshTokenID(ctx, refreshToken.ID)

	// Atomic revocation
	err = uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.RefreshTokenRepo.Revoke(txCtx, refreshToken.ID, "User logout"); err != nil {
			return err
		}
		if session != nil && session.IsActive() {
			if err := uc.UserSessionRepo.Revoke(txCtx, session.ID); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return errors.ErrInternal("failed to revoke session").WithError(err)
	}

	// Blacklist access token (outside transaction, fire-and-forget for Redis).
	// Use context.WithoutCancel so blacklist completes even if HTTP client disconnects.
	if req.AccessTokenJTI != "" {
		ttl := time.Until(req.AccessTokenExp)
		if ttl > 0 {
			_ = uc.TokenBlacklistStore.BlacklistToken(context.WithoutCancel(ctx), req.AccessTokenJTI, ttl)
		}
	}

	return nil
}

package internal

import (
	"context"
	"fmt"
	"time"

	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"
)

func (uc *usecase) LogoutAll(ctx context.Context, req *authdto.LogoutAllRequest) error {
	// Atomic revocation of all tokens and sessions
	err := uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.RefreshTokenRepo.RevokeAllByUserID(txCtx, req.UserID, "User logout all"); err != nil {
			return fmt.Errorf("revoke all refresh tokens: %w", err)
		}
		if err := uc.UserSessionRepo.RevokeAllByUserID(txCtx, req.UserID); err != nil {
			return fmt.Errorf("revoke all sessions: %w", err)
		}
		return nil
	})
	if err != nil {
		return errors.ErrInternal("failed to revoke all sessions").WithError(err)
	}

	// User-level blacklist: all tokens issued before now are invalid
	// TTL = access token expiry (15 min from config)
	ttl := uc.Config.JWT.AccessExpiry
	if ttl <= 0 {
		ttl = 15 * time.Minute // fallback
	}
	// Use context.WithoutCancel so blacklist completes even if HTTP client disconnects.
	_ = uc.TokenBlacklistStore.BlacklistUser(context.WithoutCancel(ctx), req.UserID, time.Now(), ttl)

	return nil
}

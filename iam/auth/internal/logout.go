package internal

import (
	"context"

	"iam-service/pkg/errors"
)

func (uc *usecase) Logout(ctx context.Context, token string) error {
	if token == "" {
		return errors.ErrBadRequest("refresh token is required")
	}

	tokenHash := hashToken(token)

	refreshToken, err := uc.RefreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.IsNotFound(err) {
			return errors.ErrUnauthorized("Invalid or expired token")
		}
		return errors.ErrInternal("failed to verify token").WithError(err)
	}

	if refreshToken.RevokedAt != nil {
		return errors.ErrUnauthorized("Token already revoked")
	}

	if refreshToken.IsExpired() {
		return errors.ErrUnauthorized("Token has expired")
	}

	if err := uc.RefreshTokenRepo.Revoke(ctx, refreshToken.RefreshTokenID, "User logout"); err != nil {
		return errors.ErrInternal("failed to revoke token").WithError(err)
	}

	return nil
}

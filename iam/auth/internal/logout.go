package internal

import (
	"iam-service/entity"
	"iam-service/pkg/errors"
	"time"

	"gorm.io/gorm"
)

func (uc *usecase) Logout(token string) error {
	if token == "" {
		return errors.ErrBadRequest("refresh token is required")
	}

	tokenHash := hashToken(token)

	var refreshToken entity.RefreshToken
	err := uc.DB.Where("token_hash = ?", tokenHash).First(&refreshToken).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrUnauthorized("Invalid or expired token")
		}
		return errors.ErrInternal("failed to verify token").WithError(err)
	}

	if refreshToken.RevokedAt != nil {
		return errors.ErrUnauthorized("Token already revoked")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return errors.ErrUnauthorized("Token has expired")
	}

	now := time.Now()
	reason := "User logout"
	refreshToken.RevokedAt = &now
	refreshToken.RevokedReason = &reason

	if err := uc.DB.Save(&refreshToken).Error; err != nil {
		return errors.ErrInternal("failed to revoke token").WithError(err)
	}

	return nil
}

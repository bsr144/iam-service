package postgres

import (
	"context"
	"time"

	"iam-service/entity"
	"iam-service/iam/auth/contract"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) contract.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *entity.RefreshToken) error {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return errors.TranslatePostgres(err)
	}
	return nil
}

func (r *refreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	var token entity.RefreshToken
	err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		return nil, errors.TranslatePostgres(err)
	}
	return &token, nil
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID, reason string) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Model(&entity.RefreshToken{}).
		Where("refresh_token_id = ?", id).
		Updates(map[string]interface{}{
			"revoked_at":     now,
			"revoked_reason": reason,
		}).Error; err != nil {
		return errors.TranslatePostgres(err)
	}
	return nil
}

func (r *refreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID, reason string) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Model(&entity.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Updates(map[string]interface{}{
			"revoked_at":     now,
			"revoked_reason": reason,
		}).Error; err != nil {
		return errors.TranslatePostgres(err)
	}
	return nil
}

func (r *refreshTokenRepository) RevokeByFamily(ctx context.Context, tokenFamily uuid.UUID, reason string) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Model(&entity.RefreshToken{}).
		Where("token_family = ? AND revoked_at IS NULL", tokenFamily).
		Updates(map[string]interface{}{
			"revoked_at":     now,
			"revoked_reason": reason,
		}).Error; err != nil {
		return errors.TranslatePostgres(err)
	}
	return nil
}

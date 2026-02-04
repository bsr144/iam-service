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

type pinVerificationLogRepository struct {
	db *gorm.DB
}

func NewPINVerificationLogRepository(db *gorm.DB) contract.PINVerificationLogRepository {
	return &pinVerificationLogRepository{db: db}
}

func (r *pinVerificationLogRepository) Create(ctx context.Context, log *entity.PINVerificationLog) error {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return errors.TranslatePostgres(err)
	}
	return nil
}

func (r *pinVerificationLogRepository) CountRecentFailures(ctx context.Context, userID uuid.UUID, since int) (int, error) {
	var count int64
	sinceTime := time.Now().Add(-time.Duration(since) * time.Minute)
	err := r.db.WithContext(ctx).
		Model(&entity.PINVerificationLog{}).
		Where("user_id = ? AND result = false AND created_at > ?", userID, sinceTime).
		Count(&count).Error
	if err != nil {
		return 0, errors.TranslatePostgres(err)
	}
	return int(count), nil
}

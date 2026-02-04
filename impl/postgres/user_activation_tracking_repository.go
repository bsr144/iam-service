package postgres

import (
	"context"

	"iam-service/entity"
	"iam-service/iam/auth/contract"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userActivationTrackingRepository struct {
	db *gorm.DB
}

func NewUserActivationTrackingRepository(db *gorm.DB) contract.UserActivationTrackingRepository {
	return &userActivationTrackingRepository{db: db}
}

func (r *userActivationTrackingRepository) Create(ctx context.Context, tracking *entity.UserActivationTracking) error {
	if err := r.db.WithContext(ctx).Create(tracking).Error; err != nil {
		return errors.TranslatePostgres(err)
	}
	return nil
}

func (r *userActivationTrackingRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserActivationTracking, error) {
	var tracking entity.UserActivationTracking
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&tracking).Error
	if err != nil {
		return nil, errors.TranslatePostgres(err)
	}
	return &tracking, nil
}

func (r *userActivationTrackingRepository) Update(ctx context.Context, tracking *entity.UserActivationTracking) error {
	if err := r.db.WithContext(ctx).Save(tracking).Error; err != nil {
		return errors.TranslatePostgres(err)
	}
	return nil
}

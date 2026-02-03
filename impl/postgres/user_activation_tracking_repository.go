package postgres

import (
	"context"
	"errors"

	"iam-service/entity"
	"iam-service/iam/auth/contract"

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
	return r.db.WithContext(ctx).Create(tracking).Error
}

func (r *userActivationTrackingRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserActivationTracking, error) {
	var tracking entity.UserActivationTracking
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&tracking).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tracking, nil
}

func (r *userActivationTrackingRepository) Update(ctx context.Context, tracking *entity.UserActivationTracking) error {
	return r.db.WithContext(ctx).Save(tracking).Error
}

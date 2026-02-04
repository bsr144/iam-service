package postgres

import (
	"context"

	"iam-service/entity"
	"iam-service/iam/auth/contract"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) contract.UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) Create(ctx context.Context, profile *entity.UserProfile) error {
	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		return errors.TranslatePostgres(err)
	}
	return nil
}

func (r *userProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error) {
	var profile entity.UserProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, errors.TranslatePostgres(err)
	}
	return &profile, nil
}

func (r *userProfileRepository) Update(ctx context.Context, profile *entity.UserProfile) error {
	if err := r.db.WithContext(ctx).Save(profile).Error; err != nil {
		return errors.TranslatePostgres(err)
	}
	return nil
}

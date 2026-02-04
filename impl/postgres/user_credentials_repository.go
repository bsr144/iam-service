package postgres

import (
	"context"

	"iam-service/entity"
	"iam-service/iam/auth/contract"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userCredentialsRepository struct {
	db *gorm.DB
}

func NewUserCredentialsRepository(db *gorm.DB) contract.UserCredentialsRepository {
	return &userCredentialsRepository{db: db}
}

func (r *userCredentialsRepository) Create(ctx context.Context, credentials *entity.UserCredentials) error {
	if err := r.db.WithContext(ctx).Create(credentials).Error; err != nil {
		return errors.TranslatePostgres(err)
	}
	return nil
}

func (r *userCredentialsRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserCredentials, error) {
	var credentials entity.UserCredentials
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&credentials).Error
	if err != nil {
		return nil, errors.TranslatePostgres(err)
	}
	return &credentials, nil
}

func (r *userCredentialsRepository) Update(ctx context.Context, credentials *entity.UserCredentials) error {
	if err := r.db.WithContext(ctx).Save(credentials).Error; err != nil {
		return errors.TranslatePostgres(err)
	}
	return nil
}

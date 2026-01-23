package postgres

import (
	"context"
	"errors"

	"iam-service/entity"
	"iam-service/internal/auth/contract"

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
	return r.db.WithContext(ctx).Create(credentials).Error
}

func (r *userCredentialsRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserCredentials, error) {
	var credentials entity.UserCredentials
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&credentials).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &credentials, nil
}

func (r *userCredentialsRepository) Update(ctx context.Context, credentials *entity.UserCredentials) error {
	return r.db.WithContext(ctx).Save(credentials).Error
}

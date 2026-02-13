package postgres

import (
	"context"

	"iam-service/entity"
	"iam-service/iam/auth/contract"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userCredentialsRepository struct {
	baseRepository
}

func NewUserCredentialsRepository(db *gorm.DB) contract.UserCredentialsRepository {
	return &userCredentialsRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *userCredentialsRepository) Create(ctx context.Context, credentials *entity.UserCredentials) error {
	if err := r.getDB(ctx).Create(credentials).Error; err != nil {
		return translateError(err, "user credentials")
	}
	return nil
}

func (r *userCredentialsRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserCredentials, error) {
	var credentials entity.UserCredentials
	err := r.getDB(ctx).Where("user_id = ?", userID).First(&credentials).Error
	if err != nil {
		return nil, translateError(err, "user credentials")
	}
	return &credentials, nil
}

func (r *userCredentialsRepository) Update(ctx context.Context, credentials *entity.UserCredentials) error {
	if err := r.getDB(ctx).Save(credentials).Error; err != nil {
		return translateError(err, "user credentials")
	}
	return nil
}

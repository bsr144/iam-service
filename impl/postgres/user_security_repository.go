package postgres

import (
	"context"

	"iam-service/entity"
	"iam-service/iam/auth/contract"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userSecurityRepository struct {
	db *gorm.DB
}

func NewUserSecurityRepository(db *gorm.DB) contract.UserSecurityRepository {
	return &userSecurityRepository{db: db}
}

func (r *userSecurityRepository) Create(ctx context.Context, security *entity.UserSecurity) error {
	if err := r.db.WithContext(ctx).Create(security).Error; err != nil {
		return errors.TranslatePostgres(err)
	}
	return nil
}

func (r *userSecurityRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserSecurity, error) {
	var security entity.UserSecurity
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&security).Error
	if err != nil {
		return nil, errors.TranslatePostgres(err)
	}
	return &security, nil
}

func (r *userSecurityRepository) Update(ctx context.Context, security *entity.UserSecurity) error {
	if err := r.db.WithContext(ctx).Save(security).Error; err != nil {
		return errors.TranslatePostgres(err)
	}
	return nil
}

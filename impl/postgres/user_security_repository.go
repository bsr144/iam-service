package postgres

import (
	"context"
	"errors"

	"iam-service/entity"
	"iam-service/iam/auth/contract"

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
	return r.db.WithContext(ctx).Create(security).Error
}

func (r *userSecurityRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserSecurity, error) {
	var security entity.UserSecurity
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&security).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &security, nil
}

func (r *userSecurityRepository) Update(ctx context.Context, security *entity.UserSecurity) error {
	return r.db.WithContext(ctx).Save(security).Error
}

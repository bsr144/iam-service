package postgres

import (
	"context"
	"time"

	"iam-service/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRoleRepository struct {
	baseRepository
}

func NewUserRoleRepository(db *gorm.DB) *userRoleRepository {
	return &userRoleRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *userRoleRepository) Create(ctx context.Context, userRole *entity.UserRole) error {
	if err := r.getDB(ctx).Create(userRole).Error; err != nil {
		return translateError(err, "user role")
	}
	return nil
}

func (r *userRoleRepository) ListActiveByUserID(ctx context.Context, userID uuid.UUID, productID *uuid.UUID) ([]entity.UserRole, error) {
	var userRoles []entity.UserRole
	now := time.Now()

	query := r.getDB(ctx).Where("user_id = ? AND deleted_at IS NULL", userID).
		Where("effective_from <= ?", now).
		Where("effective_to IS NULL OR effective_to > ?", now)

	if productID != nil {
		query = query.Where("product_id = ? OR product_id IS NULL", *productID)
	}

	if err := query.Find(&userRoles).Error; err != nil {
		return nil, translateError(err, "user roles")
	}
	return userRoles, nil
}

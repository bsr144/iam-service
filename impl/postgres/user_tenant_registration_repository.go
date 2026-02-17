package postgres

import (
	"context"

	"iam-service/entity"
	"iam-service/iam/auth/contract"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userTenantRegistrationRepository struct {
	baseRepository
}

func NewUserTenantRegistrationRepository(db *gorm.DB) contract.UserTenantRegistrationRepository {
	return &userTenantRegistrationRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *userTenantRegistrationRepository) ListActiveByUserID(ctx context.Context, userID uuid.UUID) ([]entity.UserTenantRegistration, error) {
	var registrations []entity.UserTenantRegistration
	err := r.getDB(ctx).
		Where("user_id = ? AND status = ? AND deleted_at IS NULL", userID, entity.UTRStatusActive).
		Find(&registrations).Error
	if err != nil {
		return nil, translateError(err, "user tenant registration")
	}
	return registrations, nil
}

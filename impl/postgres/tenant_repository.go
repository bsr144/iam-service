package postgres

import (
	"context"

	"iam-service/entity"
	"iam-service/iam/auth/contract"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tenantRepository struct {
	baseRepository
}

func NewTenantRepository(db *gorm.DB) contract.TenantRepository {
	return &tenantRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *tenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Tenant, error) {
	var tenant entity.Tenant
	err := r.getDB(ctx).Where("id = ?", id).First(&tenant).Error
	if err != nil {
		return nil, translateError(err, "tenant")
	}
	return &tenant, nil
}

func (r *tenantRepository) GetBySlug(ctx context.Context, slug string) (*entity.Tenant, error) {
	var tenant entity.Tenant
	err := r.getDB(ctx).Where("slug = ?", slug).First(&tenant).Error
	if err != nil {
		return nil, translateError(err, "tenant")
	}
	return &tenant, nil
}

func (r *tenantRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.getDB(ctx).Model(&entity.Tenant{}).
		Where("id = ?", id).
		Count(&count).Error
	if err != nil {
		return false, translateError(err, "tenant")
	}
	return count > 0, nil
}

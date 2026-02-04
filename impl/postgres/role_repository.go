package postgres

import (
	"context"

	"iam-service/entity"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *roleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return errors.TranslatePostgres(err)
	}
	return nil
}

func (r *roleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).Where("role_id = ?", id).First(&role).Error
	if err != nil {
		return nil, errors.TranslatePostgres(err)
	}
	return &role, nil
}

func (r *roleRepository) GetByName(ctx context.Context, tenantID uuid.UUID, name string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND name = ?", tenantID, name).First(&role).Error
	if err != nil {
		return nil, errors.TranslatePostgres(err)
	}
	return &role, nil
}

func (r *roleRepository) GetByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND email = ?", tenantID, email).First(&role).Error
	if err != nil {
		return nil, errors.TranslatePostgres(err)
	}
	return &role, nil
}

func (r *roleRepository) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND code = ?", tenantID, code).First(&role).Error
	if err != nil {
		return nil, errors.TranslatePostgres(err)
	}
	return &role, nil
}

func (r *roleRepository) Update(ctx context.Context, role *entity.Role) error {
	if err := r.db.WithContext(ctx).Save(role).Error; err != nil {
		return errors.TranslatePostgres(err)
	}
	return nil
}

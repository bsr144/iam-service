package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"iam-service/entity"
	"iam-service/internal/user/contract"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("user_id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND email = ?", tenantID, email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmailAnyTenant(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) EmailExistsInTenant(ctx context.Context, tenantID uuid.UUID, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).
		Where("tenant_id = ? AND email = ?", tenantID, email).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&entity.User{}).
		Where("user_id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"is_active":  false,
			"updated_at": now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userRepository) List(ctx context.Context, filter *contract.UserListFilter) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.User{}).Where("deleted_at IS NULL")

	if filter.TenantID != nil {
		query = query.Where("tenant_id = ?", *filter.TenantID)
	}

	if filter.BranchID != nil {
		query = query.Where("branch_id = ?", *filter.BranchID)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where(
			"email ILIKE ? OR user_id IN (SELECT user_id FROM user_profiles WHERE first_name ILIKE ? OR last_name ILIKE ?)",
			searchTerm, searchTerm, searchTerm,
		)
	}

	if filter.RoleID != nil {
		query = query.Where(
			"user_id IN (SELECT user_id FROM user_roles WHERE role_id = ? AND deleted_at IS NULL AND (effective_to IS NULL OR effective_to > NOW()))",
			*filter.RoleID,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	validSortColumns := map[string]string{
		"created_at": "created_at",
		"updated_at": "updated_at",
		"email":      "email",
	}

	sortColumn := "created_at"
	if col, ok := validSortColumns[filter.SortBy]; ok {
		sortColumn = col
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	offset := (filter.Page - 1) * filter.PerPage
	err := query.
		Order(fmt.Sprintf("%s %s", sortColumn, sortOrder)).
		Offset(offset).
		Limit(filter.PerPage).
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) GetPendingApprovalUsers(ctx context.Context, tenantID uuid.UUID) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Where("user_id IN (SELECT user_id FROM user_activation_tracking WHERE awaiting_admin_approval = true AND admin_approval_at IS NULL)").
		Order("created_at DESC").
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

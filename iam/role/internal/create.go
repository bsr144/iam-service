package internal

import (
	"context"
	"iam-service/entity"
	"iam-service/iam/role/roledto"
	"iam-service/pkg/errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (uc *usecase) Create(ctx context.Context, req *roledto.CreateRequest) (*roledto.CreateResponse, error) {
	tenantExists, err := uc.TenantRepo.Exists(ctx, req.TenantID)
	if err != nil {
		return nil, errors.ErrInternal("failed to verify tenant").WithError(err)
	}
	if !tenantExists {
		return nil, errors.ErrTenantNotFound()
	}

	existingRole, err := uc.RoleRepo.GetByCode(ctx, req.TenantID, req.Code)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.ErrInternal("failed to check role existence").WithError(err)
	}
	if existingRole != nil {
		return nil, errors.ErrConflict("Role with this code already exists in the tenant")
	}

	scopeLevel := entity.ScopeLevel(req.ScopeLevel)
	if scopeLevel != entity.ScopeLevelSystem && scopeLevel != entity.ScopeLevelTenant &&
		scopeLevel != entity.ScopeLevelBranch && scopeLevel != entity.ScopeLevelSelf {
		return nil, errors.ErrValidation("invalid scope level")
	}

	now := time.Now()
	roleID, err := uuid.NewV7()
	if err != nil {
		return nil, errors.ErrInternal("failed to generate role ID").WithError(err)
	}

	role := &entity.Role{
		RoleID:      roleID,
		TenantID:    &req.TenantID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		ScopeLevel:  scopeLevel,
		IsSystem:    req.IsSystem,
		IsActive:    true,
	}

	err = uc.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return err
		}

		if len(req.Permissions) > 0 {
			for _, permissionID := range req.Permissions {
				rolePermissionID, err := uuid.NewV7()
				if err != nil {
					return err
				}

				rolePermission := &entity.RolePermission{
					RolePermissionID: rolePermissionID,
					RoleID:           roleID,
					PermissionID:     permissionID,
					CreatedAt:        now,
				}

				if err := tx.Create(rolePermission).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.ErrInternal("failed to create role").WithError(err)
	}

	response := &roledto.CreateResponse{
		RoleID:      role.RoleID,
		TenantID:    req.TenantID,
		Code:        role.Code,
		Name:        role.Name,
		Description: role.Description,
		ScopeLevel:  string(role.ScopeLevel),
		IsSystem:    role.IsSystem,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt,
	}

	return response, nil
}

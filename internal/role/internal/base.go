package internal

import (
	"iam-service/config"
	"iam-service/internal/role/contract"

	"gorm.io/gorm"
)

type usecase struct {
	DB         *gorm.DB
	Config     *config.Config
	TenantRepo contract.TenantRepository
	RoleRepo   contract.RoleRepository
}

func NewUsecase(
	db *gorm.DB,
	cfg *config.Config,
	tenantRepo contract.TenantRepository,
	roleRepo contract.RoleRepository,
) *usecase {
	return &usecase{
		DB:         db,
		Config:     cfg,
		TenantRepo: tenantRepo,
		RoleRepo:   roleRepo,
	}
}

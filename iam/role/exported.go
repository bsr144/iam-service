package role

import (
	"iam-service/config"
	"iam-service/iam/auth/contract"
	rolecontract "iam-service/iam/role/contract"
	"iam-service/iam/role/internal"

	"gorm.io/gorm"
)

type Usecase = rolecontract.Usecase

func NewUsecase(
	db *gorm.DB,
	cfg *config.Config,
	tenantRepo contract.TenantRepository,
	roleRepo contract.RoleRepository,
) Usecase {
	return internal.NewUsecase(
		db,
		cfg,
		tenantRepo,
		roleRepo,
	)
}

package user

import (
	"context"
	"iam-service/config"
	"iam-service/internal/auth/contract"
	"iam-service/internal/user/internal"
	"iam-service/internal/user/userdto"

	"gorm.io/gorm"
)

type Usecase interface {
	Create(ctx context.Context, req *userdto.CreateRequest) (*userdto.CreateResponse, error)
}

func NewUsecase(
	db *gorm.DB,
	cfg *config.Config,
	userRepo contract.UserRepository,
	userProfileRepo contract.UserProfileRepository,
	userCredentialsRepo contract.UserCredentialsRepository,
	userSecurityRepo contract.UserSecurityRepository,
	tenantRepo contract.TenantRepository,
	roleRepo contract.RoleRepository,
	userActivationTrackingRepo contract.UserActivationTrackingRepository,
) Usecase {
	return internal.NewUsecase(
		db,
		cfg,
		userRepo,
		userProfileRepo,
		userCredentialsRepo,
		userSecurityRepo,
		tenantRepo,
		roleRepo,
		userActivationTrackingRepo,
	)
}

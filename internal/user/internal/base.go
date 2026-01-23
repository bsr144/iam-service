package internal

import (
	"iam-service/config"
	"iam-service/internal/user/contract"

	"gorm.io/gorm"
)

type usecase struct {
	DB                         *gorm.DB
	Config                     *config.Config
	UserRepo                   contract.UserRepository
	UserProfileRepo            contract.UserProfileRepository
	UserCredentialsRepo        contract.UserCredentialsRepository
	UserSecurityRepo           contract.UserSecurityRepository
	TenantRepo                 contract.TenantRepository
	RoleRepo                   contract.RoleRepository
	UserActivationTrackingRepo contract.UserActivationTrackingRepository
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
) *usecase {
	return &usecase{
		DB:                         db,
		Config:                     cfg,
		UserRepo:                   userRepo,
		UserProfileRepo:            userProfileRepo,
		UserCredentialsRepo:        userCredentialsRepo,
		UserSecurityRepo:           userSecurityRepo,
		TenantRepo:                 tenantRepo,
		RoleRepo:                   roleRepo,
		UserActivationTrackingRepo: userActivationTrackingRepo,
	}
}

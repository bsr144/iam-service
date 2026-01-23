package internal

import (
	"iam-service/config"
	"iam-service/internal/auth/contract"

	"gorm.io/gorm"
)

type usecase struct {
	DB                         *gorm.DB
	Config                     *config.Config
	UserRepo                   contract.UserRepository
	UserProfileRepo            contract.UserProfileRepository
	UserCredentialsRepo        contract.UserCredentialsRepository
	UserSecurityRepo           contract.UserSecurityRepository
	EmailVerificationRepo      contract.EmailVerificationRepository
	TenantRepo                 contract.TenantRepository
	UserActivationTrackingRepo contract.UserActivationTrackingRepository
	RoleRepo                   contract.RoleRepository
	EmailService               contract.EmailService
}

func NewUsecase(
	db *gorm.DB,
	cfg *config.Config,
	userRepo contract.UserRepository,
	userProfileRepo contract.UserProfileRepository,
	userCredentialsRepo contract.UserCredentialsRepository,
	userSecurityRepo contract.UserSecurityRepository,
	emailVerificationRepo contract.EmailVerificationRepository,
	tenantRepo contract.TenantRepository,
	userActivationTrackingRepo contract.UserActivationTrackingRepository,
	roleRepository contract.RoleRepository,
	emailService contract.EmailService,
) *usecase {
	return &usecase{
		DB:                         db,
		Config:                     cfg,
		UserRepo:                   userRepo,
		UserProfileRepo:            userProfileRepo,
		UserCredentialsRepo:        userCredentialsRepo,
		UserSecurityRepo:           userSecurityRepo,
		EmailVerificationRepo:      emailVerificationRepo,
		TenantRepo:                 tenantRepo,
		UserActivationTrackingRepo: userActivationTrackingRepo,
		RoleRepo:                   roleRepository,
		EmailService:               emailService,
	}
}

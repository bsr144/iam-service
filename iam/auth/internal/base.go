package internal

import (
	"iam-service/config"
	"iam-service/iam/auth/contract"
	"iam-service/pkg/logger"
)

type usecase struct {
	TxManager                  contract.TransactionManager
	Config                     *config.Config
	UserRepo                   contract.UserRepository
	UserProfileRepo            contract.UserProfileRepository
	UserCredentialsRepo        contract.UserCredentialsRepository
	UserSecurityRepo           contract.UserSecurityRepository
	EmailVerificationRepo      contract.EmailVerificationRepository
	TenantRepo                 contract.TenantRepository
	UserActivationTrackingRepo contract.UserActivationTrackingRepository
	RoleRepo                   contract.RoleRepository
	RefreshTokenRepo           contract.RefreshTokenRepository
	UserRoleRepo               contract.UserRoleRepository
	ProductRepo                contract.ProductRepository
	PermissionRepo             contract.PermissionRepository
	EmailService               contract.EmailService
	Redis                      contract.RegistrationSessionStore
	AuditLogger                logger.AuditLogger
}

func NewUsecase(
	txManager contract.TransactionManager,
	cfg *config.Config,
	userRepo contract.UserRepository,
	userProfileRepo contract.UserProfileRepository,
	userCredentialsRepo contract.UserCredentialsRepository,
	userSecurityRepo contract.UserSecurityRepository,
	emailVerificationRepo contract.EmailVerificationRepository,
	tenantRepo contract.TenantRepository,
	userActivationTrackingRepo contract.UserActivationTrackingRepository,
	roleRepository contract.RoleRepository,
	refreshTokenRepo contract.RefreshTokenRepository,
	userRoleRepo contract.UserRoleRepository,
	productRepo contract.ProductRepository,
	permissionRepo contract.PermissionRepository,
	emailService contract.EmailService,
	redis contract.RegistrationSessionStore,
	auditLogger logger.AuditLogger,
) *usecase {
	return &usecase{
		TxManager:                  txManager,
		Config:                     cfg,
		UserRepo:                   userRepo,
		UserProfileRepo:            userProfileRepo,
		UserCredentialsRepo:        userCredentialsRepo,
		UserSecurityRepo:           userSecurityRepo,
		EmailVerificationRepo:      emailVerificationRepo,
		TenantRepo:                 tenantRepo,
		UserActivationTrackingRepo: userActivationTrackingRepo,
		RoleRepo:                   roleRepository,
		RefreshTokenRepo:           refreshTokenRepo,
		UserRoleRepo:               userRoleRepo,
		ProductRepo:                productRepo,
		PermissionRepo:             permissionRepo,
		EmailService:               emailService,
		Redis:                      redis,
		AuditLogger:                auditLogger,
	}
}

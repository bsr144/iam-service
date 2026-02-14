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
	TenantRepo                 contract.TenantRepository
	UserActivationTrackingRepo contract.UserActivationTrackingRepository
	RoleRepo                   contract.RoleRepository
	RefreshTokenRepo           contract.RefreshTokenRepository
	UserRoleRepo               contract.UserRoleRepository
	ProductRepo                contract.ProductRepository
	PermissionRepo             contract.PermissionRepository
	EmailService               contract.EmailService
	Redis                      contract.RegistrationSessionStore
	LoginRedis                 contract.LoginSessionStore
	UserSessionRepo            contract.UserSessionRepository
	UserTenantRegRepo          contract.UserTenantRegistrationRepository
	ProductsByTenantRepo       contract.ProductsByTenantRepository
	AuditLogger                logger.AuditLogger
}

func NewUsecase(
	txManager contract.TransactionManager,
	cfg *config.Config,
	userRepo contract.UserRepository,
	userProfileRepo contract.UserProfileRepository,
	userCredentialsRepo contract.UserCredentialsRepository,
	userSecurityRepo contract.UserSecurityRepository,
	tenantRepo contract.TenantRepository,
	userActivationTrackingRepo contract.UserActivationTrackingRepository,
	roleRepository contract.RoleRepository,
	refreshTokenRepo contract.RefreshTokenRepository,
	userRoleRepo contract.UserRoleRepository,
	productRepo contract.ProductRepository,
	permissionRepo contract.PermissionRepository,
	emailService contract.EmailService,
	redis contract.RegistrationSessionStore,
	loginRedis contract.LoginSessionStore,
	userSessionRepo contract.UserSessionRepository,
	userTenantRegRepo contract.UserTenantRegistrationRepository,
	productsByTenantRepo contract.ProductsByTenantRepository,
	auditLogger logger.AuditLogger,
) *usecase {
	return &usecase{
		TxManager:                  txManager,
		Config:                     cfg,
		UserRepo:                   userRepo,
		UserProfileRepo:            userProfileRepo,
		UserCredentialsRepo:        userCredentialsRepo,
		UserSecurityRepo:           userSecurityRepo,
		TenantRepo:                 tenantRepo,
		UserActivationTrackingRepo: userActivationTrackingRepo,
		RoleRepo:                   roleRepository,
		RefreshTokenRepo:           refreshTokenRepo,
		UserRoleRepo:               userRoleRepo,
		ProductRepo:                productRepo,
		PermissionRepo:             permissionRepo,
		EmailService:               emailService,
		Redis:                      redis,
		LoginRedis:                 loginRedis,
		UserSessionRepo:            userSessionRepo,
		UserTenantRegRepo:          userTenantRegRepo,
		ProductsByTenantRepo:       productsByTenantRepo,
		AuditLogger:                auditLogger,
	}
}

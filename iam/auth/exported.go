package auth

import (
	"context"
	"iam-service/config"
	"iam-service/iam/auth/authdto"
	"iam-service/iam/auth/contract"
	"iam-service/iam/auth/internal"
	"iam-service/pkg/logger"

	"github.com/google/uuid"
)

type Usecase interface {
	Logout(ctx context.Context, token string) error

	InitiateRegistration(ctx context.Context, req *authdto.InitiateRegistrationRequest) (*authdto.InitiateRegistrationResponse, error)
	VerifyRegistrationOTP(ctx context.Context, req *authdto.VerifyRegistrationOTPRequest) (*authdto.VerifyRegistrationOTPResponse, error)
	ResendRegistrationOTP(ctx context.Context, req *authdto.ResendRegistrationOTPRequest) (*authdto.ResendRegistrationOTPResponse, error)
	CompleteRegistration(ctx context.Context, req *authdto.CompleteRegistrationRequest) (*authdto.CompleteRegistrationResponse, error)
	SetPassword(ctx context.Context, req *authdto.SetPasswordRequest) (*authdto.SetPasswordResponse, error)
	CompleteProfileRegistration(ctx context.Context, req *authdto.CompleteProfileRegistrationRequest) (*authdto.CompleteProfileRegistrationResponse, error)
	GetRegistrationStatus(ctx context.Context, registrationID uuid.UUID, email string) (*authdto.RegistrationStatusResponse, error)

	InitiateLogin(ctx context.Context, req *authdto.InitiateLoginRequest) (*authdto.UnifiedLoginResponse, error)
	VerifyLoginOTP(ctx context.Context, req *authdto.VerifyLoginOTPRequest) (*authdto.VerifyLoginOTPResponse, error)
	ResendLoginOTP(ctx context.Context, req *authdto.ResendLoginOTPRequest) (*authdto.ResendLoginOTPResponse, error)
	GetLoginStatus(ctx context.Context, req *authdto.GetLoginStatusRequest) (*authdto.LoginStatusResponse, error)
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
	roleRepo contract.RoleRepository,
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
) Usecase {
	return internal.NewUsecase(
		txManager,
		cfg,
		userRepo,
		userProfileRepo,
		userCredentialsRepo,
		userSecurityRepo,
		tenantRepo,
		userActivationTrackingRepo,
		roleRepo,
		refreshTokenRepo,
		userRoleRepo,
		productRepo,
		permissionRepo,
		emailService,
		redis,
		loginRedis,
		userSessionRepo,
		userTenantRegRepo,
		productsByTenantRepo,
		auditLogger,
	)
}

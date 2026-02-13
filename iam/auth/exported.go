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

	InitiateRegistration(ctx context.Context, req *authdto.InitiateRegistrationRequest, ipAddress, userAgent string) (*authdto.InitiateRegistrationResponse, error)
	VerifyRegistrationOTP(ctx context.Context, registrationID uuid.UUID, req *authdto.VerifyRegistrationOTPRequest) (*authdto.VerifyRegistrationOTPResponse, error)
	ResendRegistrationOTP(ctx context.Context, registrationID uuid.UUID, req *authdto.ResendRegistrationOTPRequest) (*authdto.ResendRegistrationOTPResponse, error)
	CompleteRegistration(ctx context.Context, registrationID uuid.UUID, registrationToken string, req *authdto.CompleteRegistrationRequest, ipAddress, userAgent string) (*authdto.CompleteRegistrationResponse, error)
	SetPassword(ctx context.Context, registrationID uuid.UUID, registrationToken string, req *authdto.SetPasswordRequest) (*authdto.SetPasswordResponse, error)
	CompleteProfileRegistration(ctx context.Context, registrationID uuid.UUID, registrationToken string, req *authdto.CompleteProfileRegistrationRequest) (*authdto.CompleteProfileRegistrationResponse, error)
	GetRegistrationStatus(ctx context.Context, registrationID uuid.UUID, email string) (*authdto.RegistrationStatusResponse, error)
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
		auditLogger,
	)
}

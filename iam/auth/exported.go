package auth

import (
	"context"
	"iam-service/config"
	"iam-service/iam/auth/authdto"
	"iam-service/iam/auth/contract"
	"iam-service/iam/auth/internal"

	"github.com/google/uuid"
)

type Usecase interface {
	Login(ctx context.Context, req *authdto.LoginRequest) (*authdto.LoginResponse, error)
	Logout(ctx context.Context, token string) error
	CompleteProfile(ctx context.Context, req *authdto.CompleteProfileRequest) (*authdto.CompleteProfileResponse, error)
	CreatePIN(ctx context.Context, userID string, newPIN string) error
	SetupPIN(ctx context.Context, userID uuid.UUID, req *authdto.SetupPINRequest) (*authdto.SetupPINResponse, error)
	VerifyOTP(ctx context.Context, req *authdto.VerifyOTPRequest) (*authdto.VerifyOTPResponse, error)
	Register(ctx context.Context, req *authdto.RegisterRequest) (*authdto.RegisterResponse, error)
	RegisterSpecialAccount(ctx context.Context, req *authdto.RegisterSpecialAccountRequest) (*authdto.RegisterSpecialAccountResponse, error)
	ResendOTP(ctx context.Context, req *authdto.ResendOTPRequest) (*authdto.ResendOTPResponse, error)
	RequestPasswordReset(ctx context.Context, req *authdto.RequestPasswordResetRequest) (*authdto.RequestPasswordResetResponse, error)
	ResetPassword(ctx context.Context, req *authdto.ResetPasswordRequest) (*authdto.ResetPasswordResponse, error)

	InitiateRegistration(ctx context.Context, tenantID uuid.UUID, req *authdto.InitiateRegistrationRequest, ipAddress, userAgent string) (*authdto.InitiateRegistrationResponse, error)
	VerifyRegistrationOTP(ctx context.Context, tenantID, registrationID uuid.UUID, req *authdto.VerifyRegistrationOTPRequest) (*authdto.VerifyRegistrationOTPResponse, error)
	ResendRegistrationOTP(ctx context.Context, tenantID, registrationID uuid.UUID, req *authdto.ResendRegistrationOTPRequest) (*authdto.ResendRegistrationOTPResponse, error)
	CompleteRegistration(ctx context.Context, tenantID, registrationID uuid.UUID, registrationToken string, req *authdto.CompleteRegistrationRequest, ipAddress, userAgent string) (*authdto.CompleteRegistrationResponse, error)
	GetRegistrationStatus(ctx context.Context, tenantID, registrationID uuid.UUID, email string) (*authdto.RegistrationStatusResponse, error)
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
	roleRepo contract.RoleRepository,
	refreshTokenRepo contract.RefreshTokenRepository,
	userRoleRepo contract.UserRoleRepository,
	productRepo contract.ProductRepository,
	permissionRepo contract.PermissionRepository,
	emailService contract.EmailService,
	redis contract.RegistrationSessionStore,
) Usecase {
	return internal.NewUsecase(
		txManager,
		cfg,
		userRepo,
		userProfileRepo,
		userCredentialsRepo,
		userSecurityRepo,
		emailVerificationRepo,
		tenantRepo,
		userActivationTrackingRepo,
		roleRepo,
		refreshTokenRepo,
		userRoleRepo,
		productRepo,
		permissionRepo,
		emailService,
		redis,
	)
}

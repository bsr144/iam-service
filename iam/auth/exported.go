package auth

import (
	"context"
	"iam-service/config"
	"iam-service/iam/auth/authdto"
	"iam-service/iam/auth/contract"
	"iam-service/iam/auth/internal"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Usecase interface {
	Login(ctx context.Context, req *authdto.LoginRequest) (*authdto.LoginResponse, error)
	Logout(token string) error
	CompleteProfile(ctx context.Context, req *authdto.CompleteProfileRequest) (*authdto.CompleteProfileResponse, error)
	CreatePIN(ctx context.Context, userID string, newPIN string) error
	SetupPIN(ctx context.Context, userID uuid.UUID, req *authdto.SetupPINRequest) (*authdto.SetupPINResponse, error)
	VerifyOTP(ctx context.Context, req *authdto.VerifyOTPRequest) (*authdto.VerifyOTPResponse, error)
	Register(ctx context.Context, req *authdto.RegisterRequest) (*authdto.RegisterResponse, error)
	RegisterSpecialAccount(ctx context.Context, req *authdto.RegisterSpecialAccountRequest) (*authdto.RegisterSpecialAccountResponse, error)
	ResendOTP(ctx context.Context, req *authdto.ResendOTPRequest) (*authdto.ResendOTPResponse, error)
	RequestPasswordReset(ctx context.Context, req *authdto.RequestPasswordResetRequest) (*authdto.RequestPasswordResetResponse, error)
	ResetPassword(ctx context.Context, req *authdto.ResetPasswordRequest) (*authdto.ResetPasswordResponse, error)
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
	roleRepo contract.RoleRepository,
	emailService contract.EmailService,
) Usecase {
	return internal.NewUsecase(
		db,
		cfg,
		userRepo,
		userProfileRepo,
		userCredentialsRepo,
		userSecurityRepo,
		emailVerificationRepo,
		tenantRepo,
		userActivationTrackingRepo,
		roleRepo,
		emailService,
	)
}

package user

import (
	"context"
	"iam-service/config"
	"iam-service/internal/user/contract"
	"iam-service/internal/user/internal"
	"iam-service/internal/user/userdto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Usecase interface {
	Create(ctx context.Context, req *userdto.CreateRequest) (*userdto.CreateResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*userdto.UserDetailResponse, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*userdto.UserDetailResponse, error)
	UpdateMe(ctx context.Context, userID uuid.UUID, req *userdto.UpdateMeRequest) (*userdto.UserDetailResponse, error)
	List(ctx context.Context, tenantID *uuid.UUID, req *userdto.ListRequest) (*userdto.ListResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *userdto.UpdateRequest) (*userdto.UserDetailResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Approve(ctx context.Context, id uuid.UUID, approverID uuid.UUID) (*userdto.ApproveResponse, error)
	Reject(ctx context.Context, id uuid.UUID, approverID uuid.UUID, req *userdto.RejectRequest) (*userdto.RejectResponse, error)
	Unlock(ctx context.Context, id uuid.UUID) (*userdto.UnlockResponse, error)
	ResetPIN(ctx context.Context, id uuid.UUID) (*userdto.ResetPINResponse, error)
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

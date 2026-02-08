package user

import (
	"context"
	"iam-service/config"
	"iam-service/iam/user/contract"
	"iam-service/iam/user/internal"
	"iam-service/iam/user/userdto"

	"github.com/google/uuid"
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
	txManager contract.TransactionManager,
	cfg *config.Config,
	userRepo contract.UserRepository,
	userProfileRepo contract.UserProfileRepository,
	userCredentialsRepo contract.UserCredentialsRepository,
	userSecurityRepo contract.UserSecurityRepository,
	tenantRepo contract.TenantRepository,
	roleRepo contract.RoleRepository,
	userActivationTrackingRepo contract.UserActivationTrackingRepository,
	userRoleRepo contract.UserRoleRepository,
) Usecase {
	return internal.NewUsecase(
		txManager,
		cfg,
		userRepo,
		userProfileRepo,
		userCredentialsRepo,
		userSecurityRepo,
		tenantRepo,
		roleRepo,
		userActivationTrackingRepo,
		userRoleRepo,
	)
}

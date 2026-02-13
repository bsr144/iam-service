package internal

import (
	"iam-service/config"
	"iam-service/entity"
	"iam-service/iam/user/contract"
	"iam-service/iam/user/userdto"
)

type usecase struct {
	TxManager                  contract.TransactionManager
	Config                     *config.Config
	UserRepo                   contract.UserRepository
	UserProfileRepo            contract.UserProfileRepository
	UserCredentialsRepo        contract.UserCredentialsRepository
	UserSecurityRepo           contract.UserSecurityRepository
	TenantRepo                 contract.TenantRepository
	RoleRepo                   contract.RoleRepository
	UserActivationTrackingRepo contract.UserActivationTrackingRepository
	UserRoleRepo               contract.UserRoleRepository
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
) *usecase {
	return &usecase{
		TxManager:                  txManager,
		Config:                     cfg,
		UserRepo:                   userRepo,
		UserProfileRepo:            userProfileRepo,
		UserCredentialsRepo:        userCredentialsRepo,
		UserSecurityRepo:           userSecurityRepo,
		TenantRepo:                 tenantRepo,
		RoleRepo:                   roleRepo,
		UserActivationTrackingRepo: userActivationTrackingRepo,
		UserRoleRepo:               userRoleRepo,
	}
}

func mapUserToDetailResponse(user *entity.User, profile *entity.UserProfile, credentials *entity.UserCredentials, security *entity.UserSecurity) *userdto.UserDetailResponse {
	resp := &userdto.UserDetailResponse{
		ID:               user.ID,
		Email:            user.Email,
		EmailVerified:    user.EmailVerified,
		IsActive:         user.IsActive,
		IsServiceAccount: user.IsServiceAccount,
		TenantID:         user.TenantID,
		BranchID:         user.BranchID,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	}

	if profile != nil {
		resp.FirstName = profile.FirstName
		resp.LastName = profile.LastName
		resp.FullName = profile.FullName()
		resp.Phone = profile.Phone
		resp.Address = profile.Address
		resp.AvatarURL = profile.AvatarURL
		resp.PreferredLanguage = profile.PreferredLanguage
		resp.Timezone = profile.Timezone
	}

	if credentials != nil {
		resp.PINSet = credentials.PINHash != nil
	}

	if security != nil {
		resp.LastLoginAt = security.LastLoginAt
	}

	return resp
}

func mapUserToListItem(user *entity.User, profile *entity.UserProfile, security *entity.UserSecurity) userdto.UserListItem {
	item := userdto.UserListItem{
		ID:            user.ID,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		IsActive:      user.IsActive,
		TenantID:      user.TenantID,
		BranchID:      user.BranchID,
		CreatedAt:     user.CreatedAt,
	}

	if profile != nil {
		item.FirstName = profile.FirstName
		item.LastName = profile.LastName
		item.FullName = profile.FullName()
		item.Phone = profile.Phone
	}

	if security != nil {
		item.LastLoginAt = security.LastLoginAt
	}

	return item
}

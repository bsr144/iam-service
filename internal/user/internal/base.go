package internal

import (
	"iam-service/config"
	"iam-service/entity"
	"iam-service/internal/user/contract"
	"iam-service/internal/user/userdto"

	"gorm.io/gorm"
)

type usecase struct {
	DB                         *gorm.DB
	Config                     *config.Config
	UserRepo                   contract.UserRepository
	UserProfileRepo            contract.UserProfileRepository
	UserCredentialsRepo        contract.UserCredentialsRepository
	UserSecurityRepo           contract.UserSecurityRepository
	TenantRepo                 contract.TenantRepository
	RoleRepo                   contract.RoleRepository
	UserActivationTrackingRepo contract.UserActivationTrackingRepository
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
) *usecase {
	return &usecase{
		DB:                         db,
		Config:                     cfg,
		UserRepo:                   userRepo,
		UserProfileRepo:            userProfileRepo,
		UserCredentialsRepo:        userCredentialsRepo,
		UserSecurityRepo:           userSecurityRepo,
		TenantRepo:                 tenantRepo,
		RoleRepo:                   roleRepo,
		UserActivationTrackingRepo: userActivationTrackingRepo,
	}
}

func mapUserToDetailResponse(user *entity.User, profile *entity.UserProfile, credentials *entity.UserCredentials, security *entity.UserSecurity) *userdto.UserDetailResponse {
	resp := &userdto.UserDetailResponse{
		ID:               user.UserID,
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
		ID:            user.UserID,
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

package internal

import (
	"context"
	"iam-service/iam/user/userdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) UpdateMe(ctx context.Context, userID uuid.UUID, req *userdto.UpdateMeRequest) (*userdto.UserDetailResponse, error) {
	user, err := uc.UserRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrInternal("failed to get user").WithError(err)
	}
	if user == nil {
		return nil, errors.ErrUserNotFound()
	}

	profile, err := uc.UserProfileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.ErrInternal("failed to get user profile").WithError(err)
	}
	if profile == nil {
		return nil, errors.ErrInternal("user profile not found")
	}

	if req.FirstName != nil {
		profile.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		profile.LastName = *req.LastName
	}
	if req.Phone != nil {
		profile.Phone = req.Phone
	}
	if req.Address != nil {
		profile.Address = req.Address
	}
	if req.PreferredLanguage != nil {
		profile.PreferredLanguage = *req.PreferredLanguage
	}
	if req.Timezone != nil {
		profile.Timezone = *req.Timezone
	}

	if err := uc.UserProfileRepo.Update(ctx, profile); err != nil {
		return nil, errors.ErrInternal("failed to update user profile").WithError(err)
	}

	credentials, err := uc.UserCredentialsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.ErrInternal("failed to get user credentials").WithError(err)
	}

	security, err := uc.UserSecurityRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.ErrInternal("failed to get user security").WithError(err)
	}

	return mapUserToDetailResponse(user, profile, credentials, security), nil
}

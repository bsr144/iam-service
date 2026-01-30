package internal

import (
	"context"
	"iam-service/internal/user/userdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) GetByID(ctx context.Context, id uuid.UUID) (*userdto.UserDetailResponse, error) {
	user, err := uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrInternal("failed to get user").WithError(err)
	}
	if user == nil {
		return nil, errors.ErrUserNotFound()
	}

	profile, err := uc.UserProfileRepo.GetByUserID(ctx, id)
	if err != nil {
		return nil, errors.ErrInternal("failed to get user profile").WithError(err)
	}

	credentials, err := uc.UserCredentialsRepo.GetByUserID(ctx, id)
	if err != nil {
		return nil, errors.ErrInternal("failed to get user credentials").WithError(err)
	}

	security, err := uc.UserSecurityRepo.GetByUserID(ctx, id)
	if err != nil {
		return nil, errors.ErrInternal("failed to get user security").WithError(err)
	}

	return mapUserToDetailResponse(user, profile, credentials, security), nil
}

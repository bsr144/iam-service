package internal

import (
	"context"

	"iam-service/iam/user/userdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) GetMe(ctx context.Context, userID uuid.UUID) (*userdto.UserDetailResponse, error) {
	user, err := uc.UserRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, err
	}

	profile, err := uc.UserProfileRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	credentials, err := uc.UserCredentialsRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	security, err := uc.UserSecurityRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	return mapUserToDetailResponse(user, profile, credentials, security), nil
}

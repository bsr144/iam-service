package internal

import (
	"context"

	"iam-service/iam/user/userdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) GetByID(ctx context.Context, callerTenantID *uuid.UUID, id uuid.UUID) (*userdto.UserDetailResponse, error) {
	user, err := uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, err
	}

	if callerTenantID != nil {
		if user.TenantID == nil || *callerTenantID != *user.TenantID {
			return nil, errors.ErrForbidden("access denied")
		}
	}

	profile, err := uc.UserProfileRepo.GetByUserID(ctx, id)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	credentials, err := uc.UserCredentialsRepo.GetByUserID(ctx, id)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	security, err := uc.UserSecurityRepo.GetByUserID(ctx, id)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	return mapUserToDetailResponse(user, profile, credentials, security), nil
}

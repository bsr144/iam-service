package internal

import (
	"context"

	"iam-service/iam/user/userdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) Update(ctx context.Context, callerTenantID *uuid.UUID, id uuid.UUID, req *userdto.UpdateRequest) (*userdto.UserDetailResponse, error) {
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

	userUpdated := false
	profileUpdated := false

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
		userUpdated = true
	}
	if req.BranchID != nil {
		user.BranchID = req.BranchID
		userUpdated = true
	}

	if userUpdated {
		if err := uc.UserRepo.Update(ctx, user); err != nil {
			return nil, errors.ErrInternal("failed to update user").WithError(err)
		}
	}

	if profile != nil {
		if req.FirstName != nil {
			profile.FirstName = *req.FirstName
			profileUpdated = true
		}
		if req.LastName != nil {
			profile.LastName = *req.LastName
			profileUpdated = true
		}
		if req.Phone != nil {
			profile.Phone = req.Phone
			profileUpdated = true
		}
		if req.Address != nil {
			profile.Address = req.Address
			profileUpdated = true
		}

		if profileUpdated {
			if err := uc.UserProfileRepo.Update(ctx, profile); err != nil {
				return nil, errors.ErrInternal("failed to update user profile").WithError(err)
			}
		}
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

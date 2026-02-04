package internal

import (
	"context"
	stderrors "errors"
	"iam-service/iam/user/userdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) Unlock(ctx context.Context, id uuid.UUID) (*userdto.UnlockResponse, error) {
	_, err := uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		if stderrors.Is(err, errors.SentinelNotFound) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, errors.ErrInternal("failed to get user").WithError(err)
	}

	security, err := uc.UserSecurityRepo.GetByUserID(ctx, id)
	if err != nil {
		if stderrors.Is(err, errors.SentinelNotFound) {
			return nil, errors.ErrInternal("user security not found")
		}
		return nil, errors.ErrInternal("failed to get user security").WithError(err)
	}

	if security.LockedUntil == nil {
		return nil, errors.ErrBadRequest("user is not locked")
	}

	security.LockedUntil = nil
	security.FailedLoginAttempts = 0

	if err := uc.UserSecurityRepo.Update(ctx, security); err != nil {
		return nil, errors.ErrInternal("failed to unlock user").WithError(err)
	}

	return &userdto.UnlockResponse{
		UserID:  id,
		Message: "User unlocked successfully",
	}, nil
}

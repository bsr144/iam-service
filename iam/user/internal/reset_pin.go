package internal

import (
	"context"

	"iam-service/iam/user/userdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) ResetPIN(ctx context.Context, id uuid.UUID) (*userdto.ResetPINResponse, error) {
	_, err := uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, err
	}

	credentials, err := uc.UserCredentialsRepo.GetByUserID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrInternal("user credentials not found")
		}
		return nil, err
	}

	if credentials.PINHash == nil {
		return nil, errors.ErrBadRequest("user does not have a PIN set")
	}

	credentials.PINHash = nil
	credentials.PINSetAt = nil
	credentials.PINChangedAt = nil
	credentials.PINExpiresAt = nil

	if err := uc.UserCredentialsRepo.Update(ctx, credentials); err != nil {
		return nil, errors.ErrInternal("failed to reset PIN").WithError(err)
	}

	return &userdto.ResetPINResponse{
		UserID:  id,
		Message: "PIN reset successfully. User will need to set a new PIN.",
	}, nil
}

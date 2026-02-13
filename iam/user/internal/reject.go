package internal

import (
	"context"

	"iam-service/iam/user/userdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) Reject(ctx context.Context, id uuid.UUID, approverID uuid.UUID, req *userdto.RejectRequest) (*userdto.RejectResponse, error) {
	user, err := uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, err
	}

	tracking, err := uc.UserActivationTrackingRepo.GetByUserID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrBadRequest("user activation tracking not found")
		}
		return nil, err
	}

	if tracking.IsActivated() {
		return nil, errors.ErrBadRequest("cannot reject an already activated user")
	}

	if tracking.IsAdminRegistered() {
		return nil, errors.ErrBadRequest("user has already been processed by admin")
	}

	err = uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := tracking.AddStatusTransition("admin_rejected: "+req.Reason, "admin"); err != nil {
			return err
		}
		if err := uc.UserActivationTrackingRepo.Update(txCtx, tracking); err != nil {
			return err
		}

		user.IsActive = false
		return uc.UserRepo.Update(txCtx, user)
	})

	if err != nil {
		return nil, errors.ErrInternal("failed to reject user").WithError(err)
	}

	return &userdto.RejectResponse{
		UserID:  id,
		Message: "User rejected successfully",
	}, nil
}

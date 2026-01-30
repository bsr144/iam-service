package internal

import (
	"context"
	"iam-service/internal/user/userdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (uc *usecase) Reject(ctx context.Context, id uuid.UUID, approverID uuid.UUID, req *userdto.RejectRequest) (*userdto.RejectResponse, error) {
	user, err := uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrInternal("failed to get user").WithError(err)
	}
	if user == nil {
		return nil, errors.ErrUserNotFound()
	}

	tracking, err := uc.UserActivationTrackingRepo.GetByUserID(ctx, id)
	if err != nil {
		return nil, errors.ErrInternal("failed to get activation tracking").WithError(err)
	}
	if tracking == nil {
		return nil, errors.ErrBadRequest("user activation tracking not found")
	}

	if tracking.IsActivated() {
		return nil, errors.ErrBadRequest("cannot reject an already activated user")
	}

	if tracking.IsAdminRegistered() {
		return nil, errors.ErrBadRequest("user has already been processed by admin")
	}

	err = uc.DB.Transaction(func(tx *gorm.DB) error {
		if err := tracking.AddStatusTransition("admin_rejected: "+req.Reason, "admin"); err != nil {
			return err
		}
		if err := tx.Save(tracking).Error; err != nil {
			return err
		}

		user.IsActive = false
		return tx.Save(user).Error
	})

	if err != nil {
		return nil, errors.ErrInternal("failed to reject user").WithError(err)
	}

	return &userdto.RejectResponse{
		UserID:  id,
		Message: "User rejected successfully",
	}, nil
}

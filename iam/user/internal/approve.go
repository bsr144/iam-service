package internal

import (
	"context"
	"time"

	"iam-service/entity"
	"iam-service/iam/user/userdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) Approve(ctx context.Context, id uuid.UUID, approverID uuid.UUID) (*userdto.ApproveResponse, error) {
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
		return nil, errors.ErrBadRequest("user is already activated")
	}

	if tracking.IsAdminRegistered() {
		return nil, errors.ErrBadRequest("user has already been approved by admin")
	}

	err = uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		now := time.Now()
		tracking.AdminCreated = true
		tracking.AdminCreatedAt = &now
		tracking.AdminCreatedBy = &approverID

		if err := tracking.AddStatusTransition("admin_approved", "admin"); err != nil {
			return err
		}

		if tracking.IsUserRegistered() {
			if err := tracking.Activate(); err != nil {
				return err
			}
			user.IsActive = true
			if err := uc.UserRepo.Update(txCtx, user); err != nil {
				return err
			}
		}

		return uc.UserActivationTrackingRepo.Update(txCtx, tracking)
	})

	if err != nil {
		return nil, errors.ErrInternal("failed to approve user").WithError(err)
	}

	return &userdto.ApproveResponse{
		UserID:  id,
		Message: "User approved successfully",
	}, nil
}

func (uc *usecase) AwaitingAdminApproval(ctx context.Context, userID uuid.UUID) (*entity.UserActivationTracking, error) {
	tracking, err := uc.UserActivationTrackingRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	if tracking.IsPendingAdminApproval() {
		return tracking, nil
	}
	return nil, nil
}

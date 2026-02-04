package internal

import (
	"context"
	stderrors "errors"
	"iam-service/entity"
	"iam-service/iam/user/userdto"
	"iam-service/pkg/errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (uc *usecase) Approve(ctx context.Context, id uuid.UUID, approverID uuid.UUID) (*userdto.ApproveResponse, error) {
	user, err := uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		if stderrors.Is(err, errors.SentinelNotFound) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, errors.ErrInternal("failed to get user").WithError(err)
	}

	tracking, err := uc.UserActivationTrackingRepo.GetByUserID(ctx, id)
	if err != nil {
		if stderrors.Is(err, errors.SentinelNotFound) {
			return nil, errors.ErrBadRequest("user activation tracking not found")
		}
		return nil, errors.ErrInternal("failed to get activation tracking").WithError(err)
	}

	if tracking.IsActivated() {
		return nil, errors.ErrBadRequest("user is already activated")
	}

	if tracking.IsAdminRegistered() {
		return nil, errors.ErrBadRequest("user has already been approved by admin")
	}

	err = uc.DB.Transaction(func(tx *gorm.DB) error {
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
			if err := tx.Save(user).Error; err != nil {
				return err
			}
		}

		return tx.Save(tracking).Error
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
		if stderrors.Is(err, errors.SentinelNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if tracking.IsPendingAdminApproval() {
		return tracking, nil
	}
	return nil, nil
}

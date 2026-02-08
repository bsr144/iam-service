package internal

import (
	"context"
	"time"

	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) CompleteProfile(ctx context.Context, req *authdto.CompleteProfileRequest) (*authdto.CompleteProfileResponse, error) {
	claims, err := uc.parseRegistrationToken(req.RegistrationToken)
	if err != nil {
		return nil, errors.ErrTokenInvalid()
	}

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return nil, errors.ErrTokenInvalid()
	}

	profile, err := uc.UserProfileRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, err
	}

	user, err := uc.UserRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, err
	}

	now := time.Now()

	if req.Address != nil {
		profile.Address = req.Address
	}
	if req.Phone != nil {
		profile.Phone = req.Phone
	}
	if req.Gender != nil {
		gender := entity.Gender(*req.Gender)
		profile.Gender = &gender
	}
	if req.MaritalStatus != nil {
		maritalStatus := entity.MaritalStatus(*req.MaritalStatus)
		profile.MaritalStatus = &maritalStatus
	}
	if req.DateOfBirth != nil {
		profile.DateOfBirth = req.DateOfBirth
	}
	if req.PlaceOfBirth != nil {
		profile.PlaceOfBirth = req.PlaceOfBirth
	}
	profile.UpdatedAt = now

	err = uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.UserProfileRepo.Update(txCtx, profile); err != nil {
			return err
		}

		tracking, err := uc.UserActivationTrackingRepo.GetByUserID(txCtx, userID)
		if err != nil {
			return err
		}
		if tracking != nil {
			if err := tracking.MarkProfileCompleted(); err != nil {
				return err
			}
			if err := tracking.MarkUserCompleted(); err != nil {
				return err
			}
			if err := uc.UserActivationTrackingRepo.Update(txCtx, tracking); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.ErrInternal("failed to complete profile").WithError(err)
	}

	return &authdto.CompleteProfileResponse{
		UserID:   userID,
		Status:   string(entity.UserStatusPendingAdminApproval),
		Email:    user.Email,
		FullName: profile.FullName(),
		Message:  "Profile completed successfully. Your registration is pending admin approval.",
	}, nil
}

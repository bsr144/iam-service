package internal

import (
	"context"
	stderrors "errors"
	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
		if stderrors.Is(err, errors.SentinelNotFound) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, errors.ErrInternal("failed to get profile").WithError(err)
	}

	user, err := uc.UserRepo.GetByID(ctx, userID)
	if err != nil {
		if stderrors.Is(err, errors.SentinelNotFound) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, errors.ErrInternal("failed to get user").WithError(err)
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

	err = uc.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(profile).Error; err != nil {
			return err
		}

		tracking, err := uc.UserActivationTrackingRepo.GetByUserID(ctx, userID)
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
			if err := tx.Save(tracking).Error; err != nil {
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

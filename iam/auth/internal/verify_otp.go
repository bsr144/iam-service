package internal

import (
	"context"
	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (uc *usecase) VerifyOTP(ctx context.Context, req *authdto.VerifyOTPRequest) (*authdto.VerifyOTPResponse, error) {
	user, err := uc.UserRepo.GetByEmail(ctx, req.TenantID, req.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to get user").WithError(err)
	}
	if user == nil {
		return nil, errors.ErrUserNotFound()
	}

	verification, err := uc.EmailVerificationRepo.GetLatestByEmail(ctx, req.Email, entity.OTPTypeRegistration)
	if err != nil {
		return nil, errors.ErrInternal("failed to get verification").WithError(err)
	}
	if verification == nil {
		return nil, errors.ErrOTPInvalid()
	}

	if verification.IsExpired() {
		return nil, errors.ErrOTPExpired()
	}

	err = bcrypt.CompareHashAndPassword([]byte(verification.OTPHash), []byte(req.OTPCode))
	if err != nil {
		return nil, errors.ErrOTPInvalid()
	}

	now := time.Now()

	err = uc.DB.Transaction(func(tx *gorm.DB) error {

		verification.VerifiedAt = &now
		if err := tx.Save(verification).Error; err != nil {
			return err
		}

		user.EmailVerified = true
		user.EmailVerifiedAt = &now
		if err := tx.Save(user).Error; err != nil {
			return err
		}

		tracking, err := uc.UserActivationTrackingRepo.GetByUserID(ctx, user.UserID)
		if err != nil {
			return err
		}
		if tracking != nil {
			if err := tracking.MarkOTPVerified(); err != nil {
				return err
			}
			if err := tx.Save(tracking).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.ErrInternal("failed to verify OTP").WithError(err)
	}

	regToken, err := uc.generateRegistrationToken(user.UserID, *user.TenantID)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate registration token").WithError(err)
	}

	return &authdto.VerifyOTPResponse{
		RegistrationToken: regToken,
		ExpiresIn:         int(time.Duration(RegistrationTokenExpiryMinutes) * time.Minute / time.Second),
		NextStep:          "complete_profile",
	}, nil
}

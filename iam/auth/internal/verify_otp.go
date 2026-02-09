package internal

import (
	"context"
	"time"

	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

func (uc *usecase) VerifyOTP(ctx context.Context, req *authdto.VerifyOTPRequest) (*authdto.VerifyOTPResponse, error) {
	user, err := uc.UserRepo.GetByEmail(ctx, req.TenantID, req.Email)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrUserNotFound()
		}
		return nil, err
	}

	verification, err := uc.EmailVerificationRepo.GetLatestByEmail(ctx, req.Email, entity.OTPTypeRegistration)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrOTPInvalid()
		}
		return nil, err
	}

	if verification.IsExpired() {
		return nil, errors.ErrOTPExpired()
	}

	err = bcrypt.CompareHashAndPassword([]byte(verification.OTPHash), []byte(req.OTPCode))
	if err != nil {
		return nil, errors.ErrOTPInvalid()
	}

	now := time.Now()

	err = uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.EmailVerificationRepo.MarkAsVerified(txCtx, verification.ID); err != nil {
			return err
		}

		user.EmailVerified = true
		user.EmailVerifiedAt = &now
		if err := uc.UserRepo.Update(txCtx, user); err != nil {
			return err
		}

		tracking, err := uc.UserActivationTrackingRepo.GetByUserID(txCtx, user.ID)
		if err != nil {
			return err
		}
		if tracking != nil {
			if err := tracking.MarkOTPVerified(); err != nil {
				return err
			}
			if err := uc.UserActivationTrackingRepo.Update(txCtx, tracking); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.ErrInternal("failed to verify OTP").WithError(err)
	}

	regToken, err := uc.generateRegistrationToken(user.ID, *user.TenantID)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate registration token").WithError(err)
	}

	return &authdto.VerifyOTPResponse{
		RegistrationToken: regToken,
		ExpiresIn:         int(time.Duration(RegistrationTokenExpiryMinutes) * time.Minute / time.Second),
		NextStep:          "complete_profile",
	}, nil
}

package internal

import (
	"context"
	"encoding/json"
	"time"

	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

func (uc *usecase) ResetPassword(ctx context.Context, req *authdto.ResetPasswordRequest) (*authdto.ResetPasswordResponse, error) {
	user, err := uc.UserRepo.GetByEmail(ctx, req.TenantID, req.Email)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrInvalidCredentials()
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.ErrUserSuspended()
	}

	verification, err := uc.EmailVerificationRepo.GetLatestByEmail(ctx, req.Email, entity.OTPTypePasswordReset)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrOTPInvalid()
		}
		return nil, err
	}

	if time.Now().After(verification.ExpiresAt) {
		return nil, errors.ErrOTPExpired()
	}

	if verification.VerifiedAt != nil {
		return nil, errors.ErrOTPInvalid()
	}

	err = bcrypt.CompareHashAndPassword([]byte(verification.OTPHash), []byte(req.OTPCode))
	if err != nil {
		return nil, errors.ErrOTPInvalid()
	}

	if err := uc.validatePassword(req.NewPassword); err != nil {
		return nil, err
	}

	credentials, err := uc.UserCredentialsRepo.GetByUserID(ctx, user.UserID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrInternal("user credentials not found")
		}
		return nil, err
	}

	if credentials.PasswordHash != nil {
		err = bcrypt.CompareHashAndPassword([]byte(*credentials.PasswordHash), []byte(req.NewPassword))
		if err == nil {
			return nil, errors.ErrValidation("new password cannot be the same as the old password")
		}
	}

	if err := uc.checkPasswordHistory(credentials, req.NewPassword); err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternal("failed to hash password").WithError(err)
	}

	err = uc.TxManager.WithTransaction(ctx, func(txCtx context.Context) error {
		now := time.Now()
		passwordHashStr := string(passwordHash)
		passwordExpiresAt := now.AddDate(0, 6, 0)

		var passwordHistory []string
		if err := json.Unmarshal(credentials.PasswordHistory, &passwordHistory); err != nil {
			passwordHistory = []string{}
		}

		if credentials.PasswordHash != nil {
			passwordHistory = append([]string{*credentials.PasswordHash}, passwordHistory...)

			if len(passwordHistory) > 5 {
				passwordHistory = passwordHistory[:5]
			}
		}

		passwordHistoryJSON, err := json.Marshal(passwordHistory)
		if err != nil {
			return errors.ErrInternal("failed to marshal password history").WithError(err)
		}

		credentials.PasswordHash = &passwordHashStr
		credentials.PasswordChangedAt = &now
		credentials.PasswordExpiresAt = &passwordExpiresAt
		credentials.PasswordHistory = passwordHistoryJSON
		credentials.UpdatedAt = now

		if err := uc.UserCredentialsRepo.Update(txCtx, credentials); err != nil {
			return errors.ErrInternal("failed to update credentials").WithError(err)
		}

		if err := uc.EmailVerificationRepo.MarkAsVerified(txCtx, verification.EmailVerificationID); err != nil {
			return errors.ErrInternal("failed to mark OTP as verified").WithError(err)
		}

		if err := uc.RefreshTokenRepo.RevokeAllByUserID(txCtx, user.UserID, "Password reset"); err != nil {
			return errors.ErrInternal("failed to revoke tokens").WithError(err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &authdto.ResetPasswordResponse{
		Message: "Password reset successful. Please login with your new password.",
	}, nil
}

func (uc *usecase) checkPasswordHistory(credentials *entity.UserCredentials, newPassword string) error {
	var passwordHistory []string
	if err := json.Unmarshal(credentials.PasswordHistory, &passwordHistory); err != nil {

		return nil
	}

	for _, oldHash := range passwordHistory {
		err := bcrypt.CompareHashAndPassword([]byte(oldHash), []byte(newPassword))
		if err == nil {
			return errors.ErrValidation("new password cannot be one of your recent passwords")
		}
	}

	return nil
}

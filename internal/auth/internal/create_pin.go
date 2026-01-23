package internal

import (
	"context"
	"encoding/json"
	"iam-service/internal/auth/authdto"
	"iam-service/pkg/errors"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc *usecase) CreatePIN(ctx context.Context, userID string, newPIN string) error {
	return errors.ErrBadRequest("use SetupPIN instead")
}

func (uc *usecase) SetupPIN(ctx context.Context, userID uuid.UUID, req *authdto.SetupPINRequest) (*authdto.SetupPINResponse, error) {
	if req.PIN != req.PINConfirm {
		return nil, errors.ErrValidation("PIN and confirmation do not match")
	}

	if err := validatePIN(req.PIN); err != nil {
		return nil, err
	}

	user, err := uc.UserRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrInternal("failed to get user").WithError(err)
	}
	if user == nil {
		return nil, errors.ErrUserNotFound()
	}

	credentials, err := uc.UserCredentialsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.ErrInternal("failed to get credentials").WithError(err)
	}
	if credentials == nil {
		return nil, errors.ErrInternal("user credentials not found")
	}

	if credentials.PINHash != nil {
		return nil, errors.ErrConflict("PIN already set. Use change PIN endpoint to update.")
	}

	pinHash, err := bcrypt.GenerateFromPassword([]byte(req.PIN), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternal("failed to hash PIN").WithError(err)
	}

	now := time.Now()
	pinExpiresAt := now.AddDate(1, 0, 0)
	pinHashStr := string(pinHash)

	credentials.PINHash = &pinHashStr
	credentials.PINSetAt = &now
	credentials.PINExpiresAt = &pinExpiresAt

	pinHistory := []string{pinHashStr}
	pinHistoryJSON, err := json.Marshal(pinHistory)
	if err != nil {
		return nil, errors.ErrInternal("failed to marshal PIN history").WithError(err)
	}
	credentials.PINHistory = pinHistoryJSON
	credentials.UpdatedAt = now

	if err := uc.UserCredentialsRepo.Update(ctx, credentials); err != nil {
		return nil, errors.ErrInternal("failed to update credentials").WithError(err)
	}

	tracking, err := uc.UserActivationTrackingRepo.GetByUserID(ctx, userID)
	if err == nil && tracking != nil {
		if err := tracking.MarkPINSet(); err == nil {
			_ = uc.UserActivationTrackingRepo.Update(ctx, tracking)
		}
	}

	response := &authdto.SetupPINResponse{
		PINSetAt:     now,
		PINExpiresAt: pinExpiresAt,
	}

	return response, nil
}
func validatePIN(pin string) error {
	if len(pin) != 6 {
		return errors.ErrValidation("PIN must be exactly 6 digits")
	}

	matched, _ := regexp.MatchString("^[0-9]{6}$", pin)
	if !matched {
		return errors.ErrValidation("PIN must contain only digits")
	}

	if pin == "123456" || pin == "654321" || pin == "012345" {
		return errors.ErrValidation("PIN cannot be a simple sequential pattern")
	}

	if pin == "000000" || pin == "111111" || pin == "222222" || pin == "333333" ||
		pin == "444444" || pin == "555555" || pin == "666666" || pin == "777777" ||
		pin == "888888" || pin == "999999" {
		return errors.ErrValidation("PIN cannot be all the same digit")
	}

	return nil
}

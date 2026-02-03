package internal

import (
	"context"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (uc *usecase) Delete(ctx context.Context, id uuid.UUID) error {
	user, err := uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		return errors.ErrInternal("failed to get user").WithError(err)
	}
	if user == nil {
		return errors.ErrUserNotFound()
	}

	if err := uc.UserRepo.Delete(ctx, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrUserNotFound()
		}
		return errors.ErrInternal("failed to delete user").WithError(err)
	}

	return nil
}

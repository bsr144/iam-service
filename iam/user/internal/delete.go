package internal

import (
	"context"
	stderrors "errors"
	"iam-service/impl/postgres"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (uc *usecase) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		if stderrors.Is(err, postgres.ErrRecordNotFound) {
			return errors.ErrUserNotFound()
		}
		return errors.ErrInternal("failed to get user").WithError(err)
	}

	if err := uc.UserRepo.Delete(ctx, id); err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return errors.ErrUserNotFound()
		}
		return errors.ErrInternal("failed to delete user").WithError(err)
	}

	return nil
}

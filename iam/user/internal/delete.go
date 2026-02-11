package internal

import (
	"context"

	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) Delete(ctx context.Context, callerTenantID *uuid.UUID, id uuid.UUID) error {
	user, err := uc.UserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.IsNotFound(err) {
			return errors.ErrUserNotFound()
		}
		return err
	}

	if callerTenantID != nil {
		if user.TenantID == nil || *callerTenantID != *user.TenantID {
			return errors.ErrForbidden("access denied")
		}
	}

	if err := uc.UserRepo.Delete(ctx, id); err != nil {
		if errors.IsNotFound(err) {
			return errors.ErrUserNotFound()
		}
		return err
	}

	return nil
}

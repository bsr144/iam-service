package internal

import (
	"context"
	"fmt"

	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) DeleteIdentity(ctx context.Context, identityID, participantID, tenantID string) error {
	iID, err := uuid.Parse(identityID)
	if err != nil {
		return errors.ErrBadRequest("invalid identity ID")
	}

	pID, err := uuid.Parse(participantID)
	if err != nil {
		return errors.ErrBadRequest("invalid participant ID")
	}

	tID, err := uuid.Parse(tenantID)
	if err != nil {
		return errors.ErrBadRequest("invalid tenant ID")
	}

	return uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		participant, err := uc.participantRepo.GetByID(txCtx, pID)
		if err != nil {
			return fmt.Errorf("get participant: %w", err)
		}

		if err := validateParticipantOwnership(participant, tID); err != nil {
			return err
		}

		if err := validateEditableState(participant); err != nil {
			return err
		}

		identity, err := uc.identityRepo.GetByID(txCtx, iID)
		if err != nil {
			return fmt.Errorf("get identity: %w", err)
		}
		if identity.ParticipantID != pID {
			return errors.ErrForbidden("identity does not belong to this participant")
		}

		if err := uc.identityRepo.SoftDelete(txCtx, iID); err != nil {
			return fmt.Errorf("soft delete identity: %w", err)
		}

		return nil
	})
}

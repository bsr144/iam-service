package internal

import (
	"context"
	"fmt"

	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) DeleteAddress(ctx context.Context, addressID, participantID, tenantID string) error {
	aID, err := uuid.Parse(addressID)
	if err != nil {
		return errors.ErrBadRequest("invalid address ID")
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

		address, err := uc.addressRepo.GetByID(txCtx, aID)
		if err != nil {
			return fmt.Errorf("get address: %w", err)
		}
		if address.ParticipantID != pID {
			return errors.ErrForbidden("address does not belong to this participant")
		}

		if err := uc.addressRepo.SoftDelete(txCtx, aID); err != nil {
			return fmt.Errorf("soft delete address: %w", err)
		}

		return nil
	})
}

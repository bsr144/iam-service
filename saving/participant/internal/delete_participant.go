package internal

import (
	"context"
	"fmt"

	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) DeleteParticipant(ctx context.Context, participantID, tenantID, userID string) error {
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

		if !participant.IsDraft() {
			return errors.ErrBadRequest("only DRAFT participants can be deleted")
		}

		if err := uc.participantRepo.SoftDelete(txCtx, pID); err != nil {
			return fmt.Errorf("soft delete participant: %w", err)
		}

		return nil
	})
}

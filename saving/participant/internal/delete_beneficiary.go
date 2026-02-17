package internal

import (
	"context"
	"fmt"

	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) DeleteBeneficiary(ctx context.Context, beneficiaryID, participantID, tenantID string) error {
	bID, err := uuid.Parse(beneficiaryID)
	if err != nil {
		return errors.ErrBadRequest("invalid beneficiary ID")
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

		beneficiary, err := uc.beneficiaryRepo.GetByID(txCtx, bID)
		if err != nil {
			return fmt.Errorf("get beneficiary: %w", err)
		}
		if beneficiary.ParticipantID != pID {
			return errors.ErrForbidden("beneficiary does not belong to this participant")
		}

		if err := uc.beneficiaryRepo.SoftDelete(txCtx, bID); err != nil {
			return fmt.Errorf("soft delete beneficiary: %w", err)
		}

		return nil
	})
}

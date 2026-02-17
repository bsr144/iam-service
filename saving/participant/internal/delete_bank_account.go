package internal

import (
	"context"
	"fmt"

	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) DeleteBankAccount(ctx context.Context, accountID, participantID, tenantID string) error {
	aID, err := uuid.Parse(accountID)
	if err != nil {
		return errors.ErrBadRequest("invalid bank account ID")
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

		account, err := uc.bankAccountRepo.GetByID(txCtx, aID)
		if err != nil {
			return fmt.Errorf("get bank account: %w", err)
		}
		if account.ParticipantID != pID {
			return errors.ErrForbidden("bank account does not belong to this participant")
		}

		if err := uc.bankAccountRepo.SoftDelete(txCtx, aID); err != nil {
			return fmt.Errorf("soft delete bank account: %w", err)
		}

		return nil
	})
}

package internal

import (
	"context"
	"fmt"

	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) DeleteFamilyMember(ctx context.Context, memberID, participantID, tenantID string) error {
	mID, err := uuid.Parse(memberID)
	if err != nil {
		return errors.ErrBadRequest("invalid family member ID")
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

		member, err := uc.familyMemberRepo.GetByID(txCtx, mID)
		if err != nil {
			return fmt.Errorf("get family member: %w", err)
		}
		if member.ParticipantID != pID {
			return errors.ErrForbidden("family member does not belong to this participant")
		}

		beneficiaries, err := uc.beneficiaryRepo.ListByParticipantID(txCtx, pID)
		if err != nil {
			return fmt.Errorf("list beneficiaries: %w", err)
		}

		for _, b := range beneficiaries {
			if b.FamilyMemberID == mID {
				return errors.ErrBadRequest("cannot delete family member that is referenced as a beneficiary")
			}
		}

		if err := uc.familyMemberRepo.SoftDelete(txCtx, mID); err != nil {
			return fmt.Errorf("soft delete family member: %w", err)
		}

		return nil
	})
}

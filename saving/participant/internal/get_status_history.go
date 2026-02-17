package internal

import (
	"context"
	"fmt"

	"iam-service/saving/participant/participantdto"

	"github.com/google/uuid"
)

func (uc *usecase) GetStatusHistory(ctx context.Context, participantID, tenantID string) ([]participantdto.StatusHistoryResponse, error) {
	pID, err := uuid.Parse(participantID)
	if err != nil {
		return nil, fmt.Errorf("invalid participant ID: %w", err)
	}

	tID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant ID: %w", err)
	}

	participant, err := uc.participantRepo.GetByID(ctx, pID)
	if err != nil {
		return nil, fmt.Errorf("get participant: %w", err)
	}

	if err := validateParticipantOwnership(participant, tID); err != nil {
		return nil, err
	}

	histories, err := uc.statusHistoryRepo.ListByParticipantID(ctx, pID)
	if err != nil {
		return nil, fmt.Errorf("list status history: %w", err)
	}

	results := make([]participantdto.StatusHistoryResponse, 0, len(histories))
	for _, h := range histories {
		results = append(results, mapStatusHistoryToResponse(h))
	}

	return results, nil
}

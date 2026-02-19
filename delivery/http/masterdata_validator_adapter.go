package rest

import (
	"context"

	authcontract "iam-service/iam/auth/contract"
	masterdatacontract "iam-service/masterdata/contract"

	"github.com/google/uuid"
)

// masterdataValidatorAdapter bridges masterdata.contract.ItemRepository to
// auth.contract.MasterdataValidator, keeping the auth domain free of masterdata imports.
type masterdataValidatorAdapter struct {
	repo masterdatacontract.ItemRepository
}

func (a *masterdataValidatorAdapter) ValidateItemCode(ctx context.Context, categoryCode, itemCode string, tenantID *uuid.UUID) (bool, error) {
	return a.repo.ValidateCode(ctx, categoryCode, itemCode, tenantID)
}

// compile-time assertion
var _ authcontract.MasterdataValidator = (*masterdataValidatorAdapter)(nil)

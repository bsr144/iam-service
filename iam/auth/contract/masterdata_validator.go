package contract

import (
	"context"

	"github.com/google/uuid"
)

// MasterdataValidator provides a boundary-crossing interface for validating
// masterdata item codes without importing the masterdata domain directly.
// The adapter implementation lives in delivery/http/fiber.go.
type MasterdataValidator interface {
	ValidateItemCode(ctx context.Context, categoryCode, itemCode string, tenantID *uuid.UUID) (bool, error)
}

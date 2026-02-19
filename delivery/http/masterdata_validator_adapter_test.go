package rest

import (
	"context"
	"errors"
	"testing"

	"iam-service/entity"
	authcontract "iam-service/iam/auth/contract"
	masterdatacontract "iam-service/masterdata/contract"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// stubItemRepository is a full implementation of masterdatacontract.ItemRepository
// that delegates only ValidateCode to the embedded mock; all other methods are no-ops.
// This keeps the test focused on the adapter's ValidateItemCode delegation.
type stubItemRepository struct {
	mock.Mock
}

// Compile-time assertion that the stub satisfies the full interface.
var _ masterdatacontract.ItemRepository = (*stubItemRepository)(nil)

func (s *stubItemRepository) ValidateCode(ctx context.Context, categoryCode, itemCode string, tenantID *uuid.UUID) (bool, error) {
	args := s.Called(ctx, categoryCode, itemCode, tenantID)
	return args.Bool(0), args.Error(1)
}

func (s *stubItemRepository) List(ctx context.Context, filter *masterdatacontract.ItemFilter) ([]*entity.MasterdataItem, int64, error) {
	return nil, 0, nil
}

func (s *stubItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.MasterdataItem, error) {
	return nil, nil
}

func (s *stubItemRepository) GetByCode(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, code string) (*entity.MasterdataItem, error) {
	return nil, nil
}

func (s *stubItemRepository) Create(ctx context.Context, item *entity.MasterdataItem) error {
	return nil
}

func (s *stubItemRepository) Update(ctx context.Context, item *entity.MasterdataItem) error {
	return nil
}

func (s *stubItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (s *stubItemRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entity.MasterdataItem, error) {
	return nil, nil
}

func (s *stubItemRepository) GetTree(ctx context.Context, categoryCode string, tenantID *uuid.UUID) ([]*entity.MasterdataItem, error) {
	return nil, nil
}

func (s *stubItemRepository) ListByParent(ctx context.Context, categoryCode, parentCode string, tenantID *uuid.UUID) ([]*entity.MasterdataItem, error) {
	return nil, nil
}

func (s *stubItemRepository) ExistsByCode(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID, code string) (bool, error) {
	return false, nil
}

func (s *stubItemRepository) GetDefaultItem(ctx context.Context, categoryID uuid.UUID, tenantID *uuid.UUID) (*entity.MasterdataItem, error) {
	return nil, nil
}

// compileTimeAdapterCheck verifies masterdataValidatorAdapter satisfies authcontract.MasterdataValidator.
var _ authcontract.MasterdataValidator = (*masterdataValidatorAdapter)(nil)

func TestMasterdataValidatorAdapter_ValidateItemCode(t *testing.T) {
	tenantID := uuid.New()

	tests := []struct {
		name         string
		categoryCode string
		itemCode     string
		tenantID     *uuid.UUID
		setup        func(*stubItemRepository)
		wantValid    bool
		wantErr      bool
	}{
		{
			name:         "delegates to repo and returns true for valid code",
			categoryCode: "GENDER",
			itemCode:     "GENDER_001",
			tenantID:     nil,
			setup: func(s *stubItemRepository) {
				s.On("ValidateCode", mock.Anything, "GENDER", "GENDER_001", (*uuid.UUID)(nil)).
					Return(true, nil)
			},
			wantValid: true,
		},
		{
			name:         "delegates to repo and returns false for unknown code",
			categoryCode: "GENDER",
			itemCode:     "GENDER_999",
			tenantID:     nil,
			setup: func(s *stubItemRepository) {
				s.On("ValidateCode", mock.Anything, "GENDER", "GENDER_999", (*uuid.UUID)(nil)).
					Return(false, nil)
			},
			wantValid: false,
		},
		{
			name:         "propagates repository error",
			categoryCode: "GENDER",
			itemCode:     "GENDER_001",
			tenantID:     nil,
			setup: func(s *stubItemRepository) {
				s.On("ValidateCode", mock.Anything, "GENDER", "GENDER_001", (*uuid.UUID)(nil)).
					Return(false, errors.New("database unavailable"))
			},
			wantErr: true,
		},
		{
			name:         "passes tenant-scoped validation with non-nil tenantID",
			categoryCode: "MARITAL_STATUS",
			itemCode:     "MARITAL_001",
			tenantID:     &tenantID,
			setup: func(s *stubItemRepository) {
				s.On("ValidateCode", mock.Anything, "MARITAL_STATUS", "MARITAL_001", &tenantID).
					Return(true, nil)
			},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoStub := new(stubItemRepository)
			tt.setup(repoStub)

			adapter := &masterdataValidatorAdapter{repo: repoStub}
			got, err := adapter.ValidateItemCode(context.Background(), tt.categoryCode, tt.itemCode, tt.tenantID)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantValid, got)
			repoStub.AssertExpectations(t)
		})
	}
}

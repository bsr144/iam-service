package internal

import (
	"context"
	"testing"

	"iam-service/entity"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUsecase_DeleteIdentity(t *testing.T) {
	tenantID := uuid.New()
	applicationID := uuid.New()
	userID := uuid.New()
	participantID := uuid.New()
	identityID := uuid.New()
	otherParticipantID := uuid.New()
	otherTenantID := uuid.New()

	tests := []struct {
		name        string
		identityID  string
		participantID string
		tenantID    string
		setup       func(*MockTransactionManager, *MockParticipantRepository, *MockParticipantIdentityRepository)
		wantErr     bool
		errKind     errors.Kind
	}{
		{
			name:        "success - deletes identity from DRAFT participant",
			identityID:  identityID.String(),
			participantID: participantID.String(),
			tenantID:    tenantID.String(),
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusDraft, tenantID, applicationID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)

				identity := createMockIdentity(participantID)
				identity.ID = identityID
				identRepo.On("GetByID", mock.Anything, identityID).Return(identity, nil)
				identRepo.On("SoftDelete", mock.Anything, identityID).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "success - deletes identity from REJECTED participant",
			identityID:  identityID.String(),
			participantID: participantID.String(),
			tenantID:    tenantID.String(),
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusRejected, tenantID, applicationID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)

				identity := createMockIdentity(participantID)
				identity.ID = identityID
				identRepo.On("GetByID", mock.Anything, identityID).Return(identity, nil)
				identRepo.On("SoftDelete", mock.Anything, identityID).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "error - invalid identity UUID",
			identityID:  "invalid-uuid",
			participantID: participantID.String(),
			tenantID:    tenantID.String(),
			setup:       func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {},
			wantErr:     true,
			errKind:     errors.KindBadRequest,
		},
		{
			name:        "error - invalid participant UUID",
			identityID:  identityID.String(),
			participantID: "invalid-uuid",
			tenantID:    tenantID.String(),
			setup:       func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {},
			wantErr:     true,
			errKind:     errors.KindBadRequest,
		},
		{
			name:        "error - invalid tenant UUID",
			identityID:  identityID.String(),
			participantID: participantID.String(),
			tenantID:    "invalid-uuid",
			setup:       func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {},
			wantErr:     true,
			errKind:     errors.KindBadRequest,
		},
		{
			name:        "error - participant not found",
			identityID:  identityID.String(),
			participantID: participantID.String(),
			tenantID:    tenantID.String(),
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				partRepo.On("GetByID", mock.Anything, participantID).Return(nil, errors.ErrNotFound("participant not found"))
			},
			wantErr: true,
			errKind: errors.KindNotFound,
		},
		{
			name:        "error - BOLA: wrong tenant",
			identityID:  identityID.String(),
			participantID: participantID.String(),
			tenantID:    otherTenantID.String(),
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusDraft, tenantID, applicationID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
			},
			wantErr: true,
			errKind: errors.KindForbidden,
		},
		{
			name:        "error - cannot delete from PENDING_APPROVAL participant",
			identityID:  identityID.String(),
			participantID: participantID.String(),
			tenantID:    tenantID.String(),
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusPendingApproval, tenantID, applicationID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
			},
			wantErr: true,
			errKind: errors.KindBadRequest,
		},
		{
			name:        "error - cannot delete from APPROVED participant",
			identityID:  identityID.String(),
			participantID: participantID.String(),
			tenantID:    tenantID.String(),
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusApproved, tenantID, applicationID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
			},
			wantErr: true,
			errKind: errors.KindBadRequest,
		},
		{
			name:        "error - identity not found",
			identityID:  identityID.String(),
			participantID: participantID.String(),
			tenantID:    tenantID.String(),
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusDraft, tenantID, applicationID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
				identRepo.On("GetByID", mock.Anything, identityID).Return(nil, errors.ErrNotFound("identity not found"))
			},
			wantErr: true,
			errKind: errors.KindNotFound,
		},
		{
			name:        "error - BOLA: identity belongs to different participant",
			identityID:  identityID.String(),
			participantID: participantID.String(),
			tenantID:    tenantID.String(),
			setup: func(txMgr *MockTransactionManager, partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository) {
				txMgr.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
				participant := createMockParticipant(entity.ParticipantStatusDraft, tenantID, applicationID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)

				identity := createMockIdentity(otherParticipantID)
				identity.ID = identityID
				identRepo.On("GetByID", mock.Anything, identityID).Return(identity, nil)
			},
			wantErr: true,
			errKind: errors.KindForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txMgr := new(MockTransactionManager)
			partRepo := new(MockParticipantRepository)
			identRepo := new(MockParticipantIdentityRepository)
			addrRepo := new(MockParticipantAddressRepository)
			bankRepo := new(MockParticipantBankAccountRepository)
			famRepo := new(MockParticipantFamilyMemberRepository)
			empRepo := new(MockParticipantEmploymentRepository)
			benRepo := new(MockParticipantBeneficiaryRepository)
			histRepo := new(MockParticipantStatusHistoryRepository)
			fileStorage := new(MockFileStorageAdapter)

			tt.setup(txMgr, partRepo, identRepo)

			uc := newTestUsecase(txMgr, partRepo, identRepo, addrRepo, bankRepo, famRepo, empRepo, benRepo, histRepo, fileStorage)

			err := uc.DeleteIdentity(context.Background(), tt.identityID, tt.participantID, tt.tenantID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errKind != 0 {
					var appErr *errors.AppError
					require.True(t, errors.As(err, &appErr))
					assert.Equal(t, tt.errKind, appErr.Kind)
				}
			} else {
				require.NoError(t, err)
			}

			txMgr.AssertExpectations(t)
			partRepo.AssertExpectations(t)
			identRepo.AssertExpectations(t)
		})
	}
}

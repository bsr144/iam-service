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

func TestUsecase_GetParticipant(t *testing.T) {
	tenantID := uuid.New()
	applicationID := uuid.New()
	userID := uuid.New()
	participantID := uuid.New()
	otherTenantID := uuid.New()

	tests := []struct {
		name          string
		participantID string
		tenantID      string
		setup         func(*MockParticipantRepository, *MockParticipantIdentityRepository, *MockParticipantAddressRepository, *MockParticipantBankAccountRepository, *MockParticipantFamilyMemberRepository, *MockParticipantEmploymentRepository, *MockParticipantBeneficiaryRepository)
		wantErr       bool
		errKind       errors.Kind
	}{
		{
			name:          "success - retrieves participant with all child entities",
			participantID: participantID.String(),
			tenantID:      tenantID.String(),
			setup: func(partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, benRepo *MockParticipantBeneficiaryRepository) {
				participant := createMockParticipant(entity.ParticipantStatusDraft, tenantID, applicationID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)

				identity := createMockIdentity(participantID)
				identRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantIdentity{identity}, nil)

				address := createMockAddress(participantID)
				addrRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantAddress{address}, nil)

				bankAccount := createMockBankAccount(participantID)
				bankRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantBankAccount{bankAccount}, nil)

				familyMember := createMockFamilyMember(participantID)
				famRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantFamilyMember{familyMember}, nil)

				employment := createMockEmployment(participantID)
				empRepo.On("GetByParticipantID", mock.Anything, participantID).Return(employment, nil)

				beneficiary := createMockBeneficiary(participantID, uuid.New())
				benRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantBeneficiary{beneficiary}, nil)
			},
			wantErr: false,
		},
		{
			name:          "success - retrieves participant without employment",
			participantID: participantID.String(),
			tenantID:      tenantID.String(),
			setup: func(partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, benRepo *MockParticipantBeneficiaryRepository) {
				participant := createMockParticipant(entity.ParticipantStatusDraft, tenantID, applicationID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)

				identRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantIdentity{}, nil)
				addrRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantAddress{}, nil)
				bankRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantBankAccount{}, nil)
				famRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantFamilyMember{}, nil)
				empRepo.On("GetByParticipantID", mock.Anything, participantID).Return(nil, errors.ErrNotFound("not found"))
				benRepo.On("ListByParticipantID", mock.Anything, participantID).Return([]*entity.ParticipantBeneficiary{}, nil)
			},
			wantErr: false,
		},
		{
			name:          "error - invalid participant UUID",
			participantID: "invalid-uuid",
			tenantID:      tenantID.String(),
			setup: func(partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, benRepo *MockParticipantBeneficiaryRepository) {
			},
			wantErr: true,
			errKind: errors.KindBadRequest,
		},
		{
			name:          "error - invalid tenant UUID",
			participantID: participantID.String(),
			tenantID:      "invalid-uuid",
			setup: func(partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, benRepo *MockParticipantBeneficiaryRepository) {
			},
			wantErr: true,
			errKind: errors.KindBadRequest,
		},
		{
			name:          "error - participant not found",
			participantID: participantID.String(),
			tenantID:      tenantID.String(),
			setup: func(partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, benRepo *MockParticipantBeneficiaryRepository) {
				partRepo.On("GetByID", mock.Anything, participantID).Return(nil, errors.ErrNotFound("participant not found"))
			},
			wantErr: true,
			errKind: errors.KindNotFound,
		},
		{
			name:          "error - BOLA: wrong tenant",
			participantID: participantID.String(),
			tenantID:      otherTenantID.String(),
			setup: func(partRepo *MockParticipantRepository, identRepo *MockParticipantIdentityRepository, addrRepo *MockParticipantAddressRepository, bankRepo *MockParticipantBankAccountRepository, famRepo *MockParticipantFamilyMemberRepository, empRepo *MockParticipantEmploymentRepository, benRepo *MockParticipantBeneficiaryRepository) {
				participant := createMockParticipant(entity.ParticipantStatusDraft, tenantID, applicationID, userID)
				participant.ID = participantID
				partRepo.On("GetByID", mock.Anything, participantID).Return(participant, nil)
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

			tt.setup(partRepo, identRepo, addrRepo, bankRepo, famRepo, empRepo, benRepo)

			uc := newTestUsecase(txMgr, partRepo, identRepo, addrRepo, bankRepo, famRepo, empRepo, benRepo, histRepo, fileStorage)

			resp, err := uc.GetParticipant(context.Background(), tt.participantID, tt.tenantID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, participantID, resp.ID)
				assert.Equal(t, tenantID, resp.TenantID)
			}

			partRepo.AssertExpectations(t)
			if !tt.wantErr {
				identRepo.AssertExpectations(t)
				addrRepo.AssertExpectations(t)
				bankRepo.AssertExpectations(t)
				famRepo.AssertExpectations(t)
				empRepo.AssertExpectations(t)
				benRepo.AssertExpectations(t)
			}
		})
	}
}

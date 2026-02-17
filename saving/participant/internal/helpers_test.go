package internal

import (
	"time"

	"iam-service/config"
	"iam-service/entity"

	"github.com/google/uuid"
)

func newTestUsecase(
	txManager *MockTransactionManager,
	participantRepo *MockParticipantRepository,
	identityRepo *MockParticipantIdentityRepository,
	addressRepo *MockParticipantAddressRepository,
	bankAccountRepo *MockParticipantBankAccountRepository,
	familyMemberRepo *MockParticipantFamilyMemberRepository,
	employmentRepo *MockParticipantEmploymentRepository,
	beneficiaryRepo *MockParticipantBeneficiaryRepository,
	statusHistoryRepo *MockParticipantStatusHistoryRepository,
	fileStorage *MockFileStorageAdapter,
) *usecase {
	return &usecase{
		cfg:               &config.Config{},
		txManager:         txManager,
		participantRepo:   participantRepo,
		identityRepo:      identityRepo,
		addressRepo:       addressRepo,
		bankAccountRepo:   bankAccountRepo,
		familyMemberRepo:  familyMemberRepo,
		employmentRepo:    employmentRepo,
		beneficiaryRepo:   beneficiaryRepo,
		statusHistoryRepo: statusHistoryRepo,
		fileStorage:       fileStorage,
	}
}

func createMockParticipant(status entity.ParticipantStatus, tenantID, applicationID, userID uuid.UUID) *entity.Participant {
	now := time.Now()
	return &entity.Participant{
		ID:            uuid.New(),
		TenantID:      tenantID,
		ApplicationID: applicationID,
		UserID:        &userID,
		FullName:      "Test Participant",
		Status:        status,
		CreatedBy:     userID,
		Version:       1,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func createMockIdentity(participantID uuid.UUID) *entity.ParticipantIdentity {
	now := time.Now()
	return &entity.ParticipantIdentity{
		ID:             uuid.New(),
		ParticipantID:  participantID,
		IdentityType:   "KTP",
		IdentityNumber: "1234567890123456",
		Version:        1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func createMockAddress(participantID uuid.UUID) *entity.ParticipantAddress {
	now := time.Now()
	countryCode := "ID"
	return &entity.ParticipantAddress{
		ID:            uuid.New(),
		ParticipantID: participantID,
		AddressType:   "DOMICILE",
		CountryCode:   &countryCode,
		IsPrimary:     true,
		Version:       1,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func createMockBankAccount(participantID uuid.UUID) *entity.ParticipantBankAccount {
	now := time.Now()
	return &entity.ParticipantBankAccount{
		ID:                uuid.New(),
		ParticipantID:     participantID,
		BankCode:          "014",
		AccountNumber:     "1234567890",
		AccountHolderName: "Test Holder",
		CurrencyCode:      "IDR",
		IsPrimary:         true,
		Version:           1,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func createMockFamilyMember(participantID uuid.UUID) *entity.ParticipantFamilyMember {
	now := time.Now()
	return &entity.ParticipantFamilyMember{
		ID:               uuid.New(),
		ParticipantID:    participantID,
		FullName:         "Test Family Member",
		RelationshipType: "SPOUSE",
		IsDependent:      true,
		Version:          1,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func createMockEmployment(participantID uuid.UUID) *entity.ParticipantEmployment {
	now := time.Now()
	hireDate := time.Now().AddDate(-1, 0, 0)
	return &entity.ParticipantEmployment{
		ID:               uuid.New(),
		ParticipantID:    participantID,
		PersonnelNumber:  strPtr("EMP001"),
		DateOfHire:       &hireDate,
		EmploymentStatus: strPtr("ACTIVE"),
		Version:          1,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func createMockBeneficiary(participantID, familyMemberID uuid.UUID) *entity.ParticipantBeneficiary {
	now := time.Now()
	return &entity.ParticipantBeneficiary{
		ID:             uuid.New(),
		ParticipantID:  participantID,
		FamilyMemberID: familyMemberID,
		AccountNumber:  strPtr("1234567890"),
		Version:        1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func createMockStatusHistory(participantID uuid.UUID, fromStatus *string, toStatus string) *entity.ParticipantStatusHistory {
	now := time.Now()
	return &entity.ParticipantStatusHistory{
		ID:            uuid.New(),
		ParticipantID: participantID,
		FromStatus:    fromStatus,
		ToStatus:      toStatus,
		ChangedBy:     uuid.New(),
		ChangedAt:     now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func strPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func uuidPtr(u uuid.UUID) *uuid.UUID {
	return &u
}

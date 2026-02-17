package internal

import (
	"iam-service/config"
	"iam-service/saving/participant/contract"
)

type usecase struct {
	cfg               *config.Config
	txManager         contract.TransactionManager
	participantRepo   contract.ParticipantRepository
	identityRepo      contract.ParticipantIdentityRepository
	addressRepo       contract.ParticipantAddressRepository
	bankAccountRepo   contract.ParticipantBankAccountRepository
	familyMemberRepo  contract.ParticipantFamilyMemberRepository
	employmentRepo    contract.ParticipantEmploymentRepository
	beneficiaryRepo   contract.ParticipantBeneficiaryRepository
	statusHistoryRepo contract.ParticipantStatusHistoryRepository
	fileStorage       contract.FileStorageAdapter
}

func NewUsecase(
	cfg *config.Config,
	txManager contract.TransactionManager,
	participantRepo contract.ParticipantRepository,
	identityRepo contract.ParticipantIdentityRepository,
	addressRepo contract.ParticipantAddressRepository,
	bankAccountRepo contract.ParticipantBankAccountRepository,
	familyMemberRepo contract.ParticipantFamilyMemberRepository,
	employmentRepo contract.ParticipantEmploymentRepository,
	beneficiaryRepo contract.ParticipantBeneficiaryRepository,
	statusHistoryRepo contract.ParticipantStatusHistoryRepository,
	fileStorage contract.FileStorageAdapter,
) contract.Usecase {
	return &usecase{
		cfg:               cfg,
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

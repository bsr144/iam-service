package participant

import (
	"iam-service/config"
	"iam-service/saving/participant/contract"
	"iam-service/saving/participant/internal"
)

type Usecase = contract.Usecase

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
) Usecase {
	return internal.NewUsecase(
		cfg,
		txManager,
		participantRepo,
		identityRepo,
		addressRepo,
		bankAccountRepo,
		familyMemberRepo,
		employmentRepo,
		beneficiaryRepo,
		statusHistoryRepo,
		fileStorage,
	)
}

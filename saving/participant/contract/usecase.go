package contract

import (
	"context"
	"io"

	"iam-service/saving/participant/participantdto"
)

type Usecase interface {
	// Participant lifecycle
	CreateParticipant(ctx context.Context, req *participantdto.CreateParticipantRequest) (*participantdto.ParticipantResponse, error)
	UpdatePersonalData(ctx context.Context, req *participantdto.UpdatePersonalDataRequest) (*participantdto.ParticipantResponse, error)
	GetParticipant(ctx context.Context, participantID, tenantID string) (*participantdto.ParticipantResponse, error)
	ListParticipants(ctx context.Context, req *participantdto.ListParticipantsRequest) (*participantdto.ListParticipantsResponse, error)
	DeleteParticipant(ctx context.Context, participantID, tenantID, userID string) error

	// Identity management
	SaveIdentity(ctx context.Context, req *participantdto.SaveIdentityRequest) (*participantdto.IdentityResponse, error)
	DeleteIdentity(ctx context.Context, identityID, participantID, tenantID string) error

	// Address management
	SaveAddress(ctx context.Context, req *participantdto.SaveAddressRequest) (*participantdto.AddressResponse, error)
	DeleteAddress(ctx context.Context, addressID, participantID, tenantID string) error

	// Bank account management
	SaveBankAccount(ctx context.Context, req *participantdto.SaveBankAccountRequest) (*participantdto.BankAccountResponse, error)
	DeleteBankAccount(ctx context.Context, accountID, participantID, tenantID string) error

	// Family member management
	SaveFamilyMember(ctx context.Context, req *participantdto.SaveFamilyMemberRequest) (*participantdto.FamilyMemberResponse, error)
	DeleteFamilyMember(ctx context.Context, memberID, participantID, tenantID string) error

	// Employment management
	SaveEmployment(ctx context.Context, req *participantdto.SaveEmploymentRequest) (*participantdto.EmploymentResponse, error)

	// Beneficiary management
	SaveBeneficiary(ctx context.Context, req *participantdto.SaveBeneficiaryRequest) (*participantdto.BeneficiaryResponse, error)
	DeleteBeneficiary(ctx context.Context, beneficiaryID, participantID, tenantID string) error

	// File management
	UploadFile(ctx context.Context, req *participantdto.UploadFileRequest, file io.Reader, fileSize int64, contentType, filename string) (*participantdto.FileUploadResponse, error)

	// Workflow actions
	SubmitParticipant(ctx context.Context, req *participantdto.SubmitParticipantRequest) (*participantdto.ParticipantResponse, error)
	ApproveParticipant(ctx context.Context, req *participantdto.ApproveParticipantRequest) (*participantdto.ParticipantResponse, error)
	RejectParticipant(ctx context.Context, req *participantdto.RejectParticipantRequest) (*participantdto.ParticipantResponse, error)

	// Status history
	GetStatusHistory(ctx context.Context, participantID, tenantID string) ([]participantdto.StatusHistoryResponse, error)
}

package internal

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"iam-service/entity"
	"iam-service/saving/participant/participantdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func validateParticipantOwnership(participant *entity.Participant, tenantID uuid.UUID) error {
	if participant.TenantID != tenantID {
		return errors.ErrForbidden("participant does not belong to this tenant")
	}
	return nil
}

func validateEditableState(participant *entity.Participant) error {
	if !participant.CanBeEdited() {
		return errors.ErrBadRequest(fmt.Sprintf("participant in %s status cannot be edited", participant.Status))
	}
	return nil
}

var allowedFieldNames = map[string]bool{
	"ktp_photo":        true,
	"passport_photo":   true,
	"family_card":      true,
	"identity_photo":   true,
	"bank_book_photo":  true,
	"supporting_doc":   true,
	"profile_photo":    true,
}

func sanitizeFieldName(fieldName string) string {
	safe := filepath.Base(fieldName)
	if safe == "." || safe == ".." || strings.ContainsAny(safe, "/\\") {
		return "unknown"
	}
	if !allowedFieldNames[safe] {
		return "unknown"
	}
	return safe
}

func sanitizeFilename(filename string) string {
	safe := filepath.Base(filename)
	if safe == "." || safe == ".." || strings.ContainsAny(safe, "/\\") {
		return "upload"
	}
	return safe
}

func generateObjectKey(tenantID, participantID uuid.UUID, fieldName, filename string) string {
	safeField := sanitizeFieldName(fieldName)
	safeFile := sanitizeFilename(filename)
	return fmt.Sprintf("participants/%s/%s/%s/%s", tenantID.String(), participantID.String(), safeField, safeFile)
}

func mapIdentityToResponse(identity *entity.ParticipantIdentity) participantdto.IdentityResponse {
	return participantdto.IdentityResponse{
		ID:                identity.ID,
		IdentityType:      identity.IdentityType,
		IdentityNumber:    identity.IdentityNumber,
		IdentityAuthority: identity.IdentityAuthority,
		IssueDate:         identity.IssueDate,
		ExpiryDate:        identity.ExpiryDate,
		PhotoFilePath:     identity.PhotoFilePath,
		Version:           identity.Version,
		CreatedAt:         identity.CreatedAt,
		UpdatedAt:         identity.UpdatedAt,
	}
}

func mapAddressToResponse(address *entity.ParticipantAddress) participantdto.AddressResponse {
	return participantdto.AddressResponse{
		ID:              address.ID,
		AddressType:     address.AddressType,
		CountryCode:     address.CountryCode,
		ProvinceCode:    address.ProvinceCode,
		CityCode:        address.CityCode,
		DistrictCode:    address.DistrictCode,
		SubdistrictCode: address.SubdistrictCode,
		PostalCode:      address.PostalCode,
		RT:              address.RT,
		RW:              address.RW,
		AddressLine:     address.AddressLine,
		IsPrimary:       address.IsPrimary,
		Version:         address.Version,
		CreatedAt:       address.CreatedAt,
		UpdatedAt:       address.UpdatedAt,
	}
}

func mapBankAccountToResponse(account *entity.ParticipantBankAccount) participantdto.BankAccountResponse {
	return participantdto.BankAccountResponse{
		ID:                account.ID,
		BankCode:          account.BankCode,
		AccountNumber:     account.AccountNumber,
		AccountHolderName: account.AccountHolderName,
		AccountType:       account.AccountType,
		CurrencyCode:      account.CurrencyCode,
		IsPrimary:         account.IsPrimary,
		IssueDate:         account.IssueDate,
		ExpiryDate:        account.ExpiryDate,
		Version:           account.Version,
		CreatedAt:         account.CreatedAt,
		UpdatedAt:         account.UpdatedAt,
	}
}

func mapFamilyMemberToResponse(member *entity.ParticipantFamilyMember) participantdto.FamilyMemberResponse {
	return participantdto.FamilyMemberResponse{
		ID:                    member.ID,
		FullName:              member.FullName,
		RelationshipType:      member.RelationshipType,
		IsDependent:           member.IsDependent,
		SupportingDocFilePath: member.SupportingDocFilePath,
		Version:               member.Version,
		CreatedAt:             member.CreatedAt,
		UpdatedAt:             member.UpdatedAt,
	}
}

func mapEmploymentToResponse(employment *entity.ParticipantEmployment) participantdto.EmploymentResponse {
	return participantdto.EmploymentResponse{
		ID:                 employment.ID,
		PersonnelNumber:    employment.PersonnelNumber,
		DateOfHire:         employment.DateOfHire,
		CorporateGroupName: employment.CorporateGroupName,
		LegalEntityCode:    employment.LegalEntityCode,
		LegalEntityName:    employment.LegalEntityName,
		BusinessUnitCode:   employment.BusinessUnitCode,
		BusinessUnitName:   employment.BusinessUnitName,
		TenantName:         employment.TenantName,
		EmploymentStatus:   employment.EmploymentStatus,
		PositionName:       employment.PositionName,
		JobLevel:           employment.JobLevel,
		LocationCode:       employment.LocationCode,
		LocationName:       employment.LocationName,
		SubLocationName:    employment.SubLocationName,
		RetirementDate:     employment.RetirementDate,
		RetirementTypeCode: employment.RetirementTypeCode,
		Version:            employment.Version,
		CreatedAt:          employment.CreatedAt,
		UpdatedAt:          employment.UpdatedAt,
	}
}

func mapBeneficiaryToResponse(beneficiary *entity.ParticipantBeneficiary) participantdto.BeneficiaryResponse {
	return participantdto.BeneficiaryResponse{
		ID:                      beneficiary.ID,
		FamilyMemberID:          beneficiary.FamilyMemberID,
		IdentityPhotoFilePath:   beneficiary.IdentityPhotoFilePath,
		FamilyCardPhotoFilePath: beneficiary.FamilyCardPhotoFilePath,
		BankBookPhotoFilePath:   beneficiary.BankBookPhotoFilePath,
		AccountNumber:           beneficiary.AccountNumber,
		Version:                 beneficiary.Version,
		CreatedAt:               beneficiary.CreatedAt,
		UpdatedAt:               beneficiary.UpdatedAt,
	}
}

func mapStatusHistoryToResponse(history *entity.ParticipantStatusHistory) participantdto.StatusHistoryResponse {
	return participantdto.StatusHistoryResponse{
		ID:         history.ID,
		FromStatus: history.FromStatus,
		ToStatus:   history.ToStatus,
		ChangedBy:  history.ChangedBy,
		Reason:     history.Reason,
		ChangedAt:  history.ChangedAt,
		CreatedAt:  history.CreatedAt,
	}
}

func (uc *usecase) buildFullParticipantResponse(ctx context.Context, participant *entity.Participant) (*participantdto.ParticipantResponse, error) {
	resp := &participantdto.ParticipantResponse{
		ID:              participant.ID,
		TenantID:        participant.TenantID,
		ApplicationID:   participant.ApplicationID,
		UserID:          participant.UserID,
		FullName:        participant.FullName,
		Gender:          participant.Gender,
		PlaceOfBirth:    participant.PlaceOfBirth,
		DateOfBirth:     participant.DateOfBirth,
		MaritalStatus:   participant.MaritalStatus,
		Citizenship:     participant.Citizenship,
		Religion:        participant.Religion,
		KTPNumber:       participant.KTPNumber,
		EmployeeNumber:  participant.EmployeeNumber,
		PhoneNumber:     participant.PhoneNumber,
		Status:          string(participant.Status),
		CreatedBy:       participant.CreatedBy,
		SubmittedBy:     participant.SubmittedBy,
		SubmittedAt:     participant.SubmittedAt,
		ApprovedBy:      participant.ApprovedBy,
		ApprovedAt:      participant.ApprovedAt,
		RejectedBy:      participant.RejectedBy,
		RejectedAt:      participant.RejectedAt,
		RejectionReason: participant.RejectionReason,
		Version:         participant.Version,
		CreatedAt:       participant.CreatedAt,
		UpdatedAt:       participant.UpdatedAt,
	}

	identities, err := uc.identityRepo.ListByParticipantID(ctx, participant.ID)
	if err != nil {
		return nil, fmt.Errorf("load identities: %w", err)
	}
	for _, identity := range identities {
		resp.Identities = append(resp.Identities, mapIdentityToResponse(identity))
	}

	addresses, err := uc.addressRepo.ListByParticipantID(ctx, participant.ID)
	if err != nil {
		return nil, fmt.Errorf("load addresses: %w", err)
	}
	for _, address := range addresses {
		resp.Addresses = append(resp.Addresses, mapAddressToResponse(address))
	}

	bankAccounts, err := uc.bankAccountRepo.ListByParticipantID(ctx, participant.ID)
	if err != nil {
		return nil, fmt.Errorf("load bank accounts: %w", err)
	}
	for _, account := range bankAccounts {
		resp.BankAccounts = append(resp.BankAccounts, mapBankAccountToResponse(account))
	}

	familyMembers, err := uc.familyMemberRepo.ListByParticipantID(ctx, participant.ID)
	if err != nil {
		return nil, fmt.Errorf("load family members: %w", err)
	}
	for _, member := range familyMembers {
		resp.FamilyMembers = append(resp.FamilyMembers, mapFamilyMemberToResponse(member))
	}

	employment, err := uc.employmentRepo.GetByParticipantID(ctx, participant.ID)
	if err != nil && !errors.IsNotFound(err) {
		return nil, fmt.Errorf("load employment: %w", err)
	}
	if employment != nil {
		empResp := mapEmploymentToResponse(employment)
		resp.Employment = &empResp
	}

	beneficiaries, err := uc.beneficiaryRepo.ListByParticipantID(ctx, participant.ID)
	if err != nil {
		return nil, fmt.Errorf("load beneficiaries: %w", err)
	}
	for _, beneficiary := range beneficiaries {
		resp.Beneficiaries = append(resp.Beneficiaries, mapBeneficiaryToResponse(beneficiary))
	}

	return resp, nil
}

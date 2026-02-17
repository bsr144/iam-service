package internal

import (
	"context"
	"encoding/json"

	"iam-service/entity"
	"iam-service/pkg/errors"
	"iam-service/saving/member/memberdto"
)

func (uc *usecase) RegisterMember(ctx context.Context, req *memberdto.RegisterRequest) (*memberdto.RegisterResponse, error) {
	_, err := uc.productRepo.GetByIDAndTenant(ctx, req.ApplicationID, req.TenantID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrNotFound("product not found in this tenant")
		}
		return nil, err
	}

	regConfig, err := uc.configRepo.GetByApplicationAndType(ctx, req.ApplicationID, "MEMBER")
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrBadRequest("member registration is not configured for this product")
		}
		return nil, err
	}

	if !regConfig.IsActive {
		return nil, errors.ErrBadRequest("member registration is currently not accepting new registrations")
	}

	existing, err := uc.utrRepo.GetByUserAndProduct(ctx, req.UserID, req.TenantID, req.ApplicationID, "MEMBER")
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if existing != nil {
		return nil, errors.ErrConflict("you already have a member registration for this product")
	}

	reg := &entity.UserTenantRegistration{
		UserID:           req.UserID,
		TenantID:         req.TenantID,
		ApplicationID:    &req.ApplicationID,
		RegistrationType: "MEMBER",
		Status:           entity.UTRStatusPendingApproval,
		Metadata:         json.RawMessage(`{}`),
	}

	if err := uc.utrRepo.Create(ctx, reg); err != nil {
		return nil, err
	}

	return &memberdto.RegisterResponse{
		ID:               reg.ID,
		Status:           string(reg.Status),
		RegistrationType: reg.RegistrationType,
		CreatedAt:        reg.CreatedAt,
	}, nil
}

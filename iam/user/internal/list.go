package internal

import (
	"context"
	"iam-service/iam/user/contract"
	"iam-service/iam/user/userdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
)

func (uc *usecase) List(ctx context.Context, tenantID *uuid.UUID, req *userdto.ListRequest) (*userdto.ListResponse, error) {
	req.SetDefaults()

	filter := &contract.UserListFilter{
		TenantID:  tenantID,
		BranchID:  req.BranchID,
		RoleID:    req.RoleID,
		Search:    req.Search,
		Page:      req.Page,
		PerPage:   req.PerPage,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	if req.Status == "active" {
		isActive := true
		filter.IsActive = &isActive
	} else if req.Status == "inactive" {
		isActive := false
		filter.IsActive = &isActive
	}

	users, total, err := uc.UserRepo.List(ctx, filter)
	if err != nil {
		return nil, errors.ErrInternal("failed to list users").WithError(err)
	}

	items := make([]userdto.UserListItem, 0, len(users))
	for _, user := range users {
		profile, _ := uc.UserProfileRepo.GetByUserID(ctx, user.ID)
		security, _ := uc.UserSecurityRepo.GetByUserID(ctx, user.ID)
		items = append(items, mapUserToListItem(user, profile, security))
	}

	totalPages := int(total) / req.PerPage
	if int(total)%req.PerPage > 0 {
		totalPages++
	}

	return &userdto.ListResponse{
		Users: items,
		Pagination: userdto.Pagination{
			Total:      total,
			Page:       req.Page,
			PerPage:    req.PerPage,
			TotalPages: totalPages,
		},
	}, nil
}

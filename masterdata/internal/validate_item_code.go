package internal

import (
	"context"

	"iam-service/masterdata/masterdatadto"
)

func (uc *usecase) ValidateItemCode(ctx context.Context, req *masterdatadto.ValidateCodeRequest) (*masterdatadto.ValidateCodeResponse, error) {
	valid, err := uc.itemRepo.ValidateCode(ctx, req.CategoryCode, req.ItemCode, req.TenantID)
	if err != nil {
		return nil, err
	}

	response := &masterdatadto.ValidateCodeResponse{
		Valid:        valid,
		CategoryCode: req.CategoryCode,
		ItemCode:     req.ItemCode,
	}

	if !valid {
		response.Message = "Item code '" + req.ItemCode + "' not found in category '" + req.CategoryCode + "'"
	}

	return response, nil
}

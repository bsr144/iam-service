package presenter

import (
	"iam-service/delivery/http/dto/response"
	"iam-service/masterdata/masterdatadto"
)

func ToItemResponse(item *masterdatadto.ItemResponse) *response.MasterdataItemResponse {
	if item == nil {
		return nil
	}
	return &response.MasterdataItemResponse{
		ID:           item.ID,
		Code:         item.Code,
		Name:         item.Name,
		AltName:      item.AltName,
		Description:  item.Description,
		ParentItemID: item.ParentItemID,
		Status:       item.Status,
		IsDefault:    item.IsDefault,
	}
}

func ToItemListResponse(items []*masterdatadto.ItemResponse) []*response.MasterdataItemResponse {
	if items == nil {
		return nil
	}
	result := make([]*response.MasterdataItemResponse, len(items))
	for i, item := range items {
		result[i] = ToItemResponse(item)
	}
	return result
}

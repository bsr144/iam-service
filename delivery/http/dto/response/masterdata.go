package response

import "github.com/google/uuid"

type MasterdataItemResponse struct {
	ID           uuid.UUID  `json:"id"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	AltName      *string    `json:"alt_name,omitempty"`
	Description  *string    `json:"description,omitempty"`
	ParentItemID *uuid.UUID `json:"parent_item_id,omitempty"`
	Status       string     `json:"status"`
	IsDefault    bool       `json:"is_default"`
}

type MasterdataCategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Status      string    `json:"status"`
}

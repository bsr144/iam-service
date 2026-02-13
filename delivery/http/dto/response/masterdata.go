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

type MasterdataItemTreeResponse struct {
	ID           uuid.UUID                     `json:"id"`
	Code         string                        `json:"code"`
	Name         string                        `json:"name"`
	AltName      *string                       `json:"alt_name,omitempty"`
	ParentItemID *uuid.UUID                    `json:"parent_item_id,omitempty"`
	Status       string                        `json:"status"`
	Children     []*MasterdataItemTreeResponse `json:"children,omitempty"`
}

type ValidateCodeResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
}

type ValidationResult struct {
	CategoryCode string `json:"category_code"`
	ItemCode     string `json:"item_code"`
	Valid        bool   `json:"valid"`
	Message      string `json:"message,omitempty"`
}

type ValidateCodesResponse struct {
	AllValid bool               `json:"all_valid"`
	Results  []ValidationResult `json:"results"`
}

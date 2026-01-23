package controller

type APIResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}

func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(errorCode string, message string) APIResponse {
	return APIResponse{
		Success: false,
		Error:   errorCode,
		Message: message,
	}
}

func ErrorResponseWithDetails(errorCode string, message string, details interface{}) APIResponse {
	return APIResponse{
		Success: false,
		Error:   errorCode,
		Message: message,
		Data:    details,
	}
}

package response

type SuccessResponse struct {
	Status     string      `json:"status"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func NewSuccessResponse(data interface{}, pagination *Pagination) SuccessResponse {
	return SuccessResponse{
		Status:     "success",
		Data:       data,
		Pagination: pagination,
	}
}

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Status:  "error",
		Message: message,
	}
}

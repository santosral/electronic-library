package handler

type SuccessResponse struct {
	Status     string      `json:"status"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type Pagination struct {
	TotalCount int `json:"total_count"`
	Limit      int `json:"limit"`
	Offset     int `json:"offset"`
}

func NewSuccessResponse(data interface{}, pagination *Pagination) SuccessResponse {
	return SuccessResponse{
		Status:     "success",
		Data:       data,
		Pagination: pagination,
	}
}

func NewErrorResponse(message string, err error) ErrorResponse {
	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
	}
	return ErrorResponse{
		Status:  "error",
		Message: message,
		Error:   errorMessage,
	}
}

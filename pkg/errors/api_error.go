package errors

import (
	"fmt"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func New(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

package types

import (
	"net/http"
)

type APIError struct {
	//  The success status of the message
	Success    bool   `json:"success" example:"false"`
	// The success or error message
	Message    string `json:"message" example:"Invalid token"`
	StatusCode int    `json:"-"`
} // @name APIResponse

func NewAPIError(success bool, err error, statusCode int) *APIError {
	return &APIError{
		Success:    success,
		Message:    err.Error(),
		StatusCode: statusCode,
	}
}

type APIFunc func(w http.ResponseWriter, r *http.Request) *APIError

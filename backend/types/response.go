package types

import "fmt"

type ApiResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"msg"`
	ErrorCode int    `json:"-"`
}

func (r ApiResponse) Error() string {
	return fmt.Sprintf("%d: %s", r.ErrorCode, r.Message)
}

func NewApiResponse(success bool, message string, httpCode int) ApiResponse {
	return ApiResponse{
		Success:   success,
		Message:   message,
		ErrorCode: httpCode,
	}
}

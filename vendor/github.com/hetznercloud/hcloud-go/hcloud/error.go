package hcloud

import "fmt"

// ErrorCode represents an error code returned from the API.
type ErrorCode string

// Error codes returned from the API.
const (
	ErrorCodeServiceError ErrorCode = "service_error" // Generic server error
	ErrorCodeLimitReached           = "limit_reached" // Rate limit reached
	ErrorCodeUnknownError           = "unknown_error" // Unknown error
)

// Error is an error returned from the API.
type Error struct {
	Code    ErrorCode
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (%s)", e.Message, e.Code)
}

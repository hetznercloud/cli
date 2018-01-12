package hcloud

import "fmt"

// ErrorCode represents an error code returned from the API.
type ErrorCode string

// Error codes returned from the API.
const (
	ErrorCodeServiceError ErrorCode = "service_error" // Generic server error
	ErrorCodeLimitReached           = "limit_reached" // Rate limit reached
	ErrorCodeUnknownError           = "unknown_error" // Unknown error
	ErrorCodeNotFound               = "not_found"     // Resource not found
	ErrorCodeInvalidInput           = "invalid_input" // Validation error
)

// Error is an error returned from the API.
type Error struct {
	Code    ErrorCode
	Message string
	Details interface{}
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (%s)", e.Message, e.Code)
}

// ErrorDetailsInvalidInput contains the details of an 'invalid_input' error.
type ErrorDetailsInvalidInput struct {
	Fields []ErrorDetailsInvalidInputField
}

// ErrorDetailsInvalidInputField contains the validation errors reported on a field.
type ErrorDetailsInvalidInputField struct {
	Name     string
	Messages []string
}

// IsError returns whether err is an API error with the given error code.
func IsError(err error, code ErrorCode) bool {
	apiErr, ok := err.(Error)
	return ok && apiErr.Code == code
}

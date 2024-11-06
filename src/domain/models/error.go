package models

import "errors"

// Common domain errors
var (
	// Repository errors
	ErrProductNotFound = errors.New("product not found")
	ErrVersionConflict = errors.New("version conflict")
	ErrInvalidProduct  = errors.New("invalid product")
	ErrLockFailed      = errors.New("failed to acquire lock")

	// API errors
	ErrInvalidRequest = errors.New("invalid request")
	ErrInternalError  = errors.New("internal server error")
)

// APIError represents an error response from the API
type APIError struct {
	Message string `json:"message"`
}

// NewAPIError creates a new API error
func NewAPIError(message string) *APIError {
	return &APIError{
		Message: message,
	}
}

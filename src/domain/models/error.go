package models

// APIError represents an API error response
type APIError struct {
	Message string `json:"error"`
}

// NewAPIError creates a new API error with the given message
func NewAPIError(message string) *APIError {
	return &APIError{
		Message: message,
	}
}

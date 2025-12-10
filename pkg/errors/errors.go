package errors

import "errors"

var (
	ErrTypeValidation = "validation"
	ErrTypeNotFound   = "not_found"
)

// AppError is a custom error type with a type field.
type AppError struct {
	Type    string
	Message string
}

// Error implements the error interface.
func (e *AppError) Error() string {
	return e.Message
}

// NewValidationError creates a validation error.
func NewValidationError(msg string) error {
	return &AppError{Type: ErrTypeValidation, Message: msg}
}

// NewNotFoundError creates a not found error.
func NewNotFoundError(msg string) error {
	return &AppError{Type: ErrTypeNotFound, Message: msg}
}

// IsValidation checks if an error is a validation error.
func IsValidation(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ErrTypeValidation
	}
	return false
}

// IsNotFound checks if an error is a not found error.
func IsNotFound(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ErrTypeNotFound
	}
	return false
}

package errors

import (
    "fmt"
    "net/http"
)

// ErrorType represents different categories of errors
type ErrorType string

const (
    ErrorTypeValidation  ErrorType = "VALIDATION_ERROR"
    ErrorTypeDatabase    ErrorType = "DATABASE_ERROR"
    ErrorTypeInternal    ErrorType = "INTERNAL_ERROR"
    ErrorTypeAuth        ErrorType = "AUTH_ERROR"
    ErrorTypeFileSystem  ErrorType = "FILESYSTEM_ERROR"
    ErrorTypeHTTPClient  ErrorType = "HTTP_CLIENT_ERROR"
)

// AppError represents an application-specific error
type AppError struct {
    Type       ErrorType         `json:"type"`
    Message    string           `json:"message"`
    Details    map[string]string `json:"details,omitempty"`
    HTTPStatus int              `json:"-"` // HTTP status code, not serialized
    err        error            `json:"-"` // Original error, not serialized
}

// Error implements the error interface
func (e *AppError) Error() string {
    if e.err != nil {
        return fmt.Sprintf("%s: %s (%s)", e.Type, e.Message, e.err.Error())
    }
    return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
    return e.err
}

// NewAppError creates a new AppError
func NewAppError(errType ErrorType, message string, err error, details map[string]string) *AppError {
    return &AppError{
        Type:       errType,
        Message:    message,
        Details:    details,
        HTTPStatus: getHTTPStatus(errType),
        err:        err,
    }
}

// Validation error helpers
func NewValidationError(message string, details map[string]string) *AppError {
    return NewAppError(ErrorTypeValidation, message, nil, details)
}

// Database error helper
func NewDatabaseError(message string, err error) *AppError {
    return NewAppError(ErrorTypeDatabase, message, err, nil)
}

// Internal error helper
func NewInternalError(message string, err error) *AppError {
    return NewAppError(ErrorTypeInternal, message, err, nil)
}

// getHTTPStatus maps error types to HTTP status codes
func getHTTPStatus(errType ErrorType) int {
    switch errType {
    case ErrorTypeValidation:
        return http.StatusBadRequest
    case ErrorTypeAuth:
        return http.StatusUnauthorized
    case ErrorTypeDatabase:
        return http.StatusInternalServerError
    case ErrorTypeFileSystem:
        return http.StatusInternalServerError
    case ErrorTypeHTTPClient:
        return http.StatusBadGateway
    default:
        return http.StatusInternalServerError
    }
}
package errors

import (
	"errors"
	"fmt"
)

// Common application errors
var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrInvalidEmail          = errors.New("invalid email address")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrMagicLinkNotFound     = errors.New("magic link not found")
	ErrMagicLinkExpired      = errors.New("magic link has expired")
	ErrMagicLinkAlreadyUsed  = errors.New("magic link has already been used")
	ErrInvalidToken          = errors.New("invalid token")
	ErrTokenExpired          = errors.New("token has expired")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrForbidden             = errors.New("forbidden")
	ErrInternalServer        = errors.New("internal server error")
	ErrEmailSendFailed       = errors.New("failed to send email")
	ErrDatabaseConnection    = errors.New("database connection error")
)

// AppError represents a custom application error with additional context
type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Error codes
const (
	CodeValidation     = "VALIDATION_ERROR"
	CodeNotFound       = "NOT_FOUND"
	CodeUnauthorized   = "UNAUTHORIZED"
	CodeForbidden      = "FORBIDDEN"
	CodeConflict       = "CONFLICT"
	CodeInternalServer = "INTERNAL_SERVER_ERROR"
	CodeBadRequest     = "BAD_REQUEST"
) 
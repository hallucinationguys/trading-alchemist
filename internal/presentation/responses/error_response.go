package responses

import (
	"net/http"

	"trading-alchemist/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents an error response
//
// swagger:model ErrorResponse
type ErrorResponse struct {
	// Error details
	Error ErrorDetail `json:"error"`
	// Success indicator (always false for errors)
	// example: false
	Success bool `json:"success"`
}

// ErrorDetail contains error information
//
// swagger:model ErrorDetail
type ErrorDetail struct {
	// Error code
	// example: VALIDATION_ERROR
	Code string `json:"code"`
	// Error message
	// example: Invalid email address
	Message string `json:"message"`
	// Optional error details
	// example: The email field must be a valid email address
	Details string `json:"details,omitempty"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code, message, details string) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
		Success: false,
	}
}

// SendError sends an error response
func SendError(c *fiber.Ctx, statusCode int, code, message string, details ...string) error {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}

	response := NewErrorResponse(code, message, detail)
	return c.Status(statusCode).JSON(response)
}

// HandleError handles application errors and sends appropriate responses
func HandleError(c *fiber.Ctx, err error) error {
	// Check if it's an AppError
	if appErr, ok := err.(*errors.AppError); ok {
		statusCode := getHTTPStatusFromErrorCode(appErr.Code)
		return SendError(c, statusCode, appErr.Code, appErr.Message, appErr.Error())
	}

	// Handle common errors
	switch err {
	case errors.ErrUserNotFound:
		return SendError(c, http.StatusNotFound, errors.CodeNotFound, "User not found")
	case errors.ErrUserAlreadyExists:
		return SendError(c, http.StatusConflict, errors.CodeConflict, "User already exists")
	case errors.ErrInvalidEmail:
		return SendError(c, http.StatusBadRequest, errors.CodeValidation, "Invalid email address")
	case errors.ErrInvalidCredentials:
		return SendError(c, http.StatusUnauthorized, errors.CodeUnauthorized, "Invalid credentials")
	case errors.ErrMagicLinkNotFound:
		return SendError(c, http.StatusNotFound, errors.CodeNotFound, "Magic link not found")
	case errors.ErrMagicLinkExpired:
		return SendError(c, http.StatusUnauthorized, errors.CodeUnauthorized, "Magic link has expired")
	case errors.ErrMagicLinkAlreadyUsed:
		return SendError(c, http.StatusUnauthorized, errors.CodeUnauthorized, "Magic link has already been used")
	case errors.ErrInvalidToken:
		return SendError(c, http.StatusUnauthorized, errors.CodeUnauthorized, "Invalid token")
	case errors.ErrTokenExpired:
		return SendError(c, http.StatusUnauthorized, errors.CodeUnauthorized, "Token has expired")
	case errors.ErrUnauthorized:
		return SendError(c, http.StatusUnauthorized, errors.CodeUnauthorized, "Unauthorized")
	case errors.ErrForbidden:
		return SendError(c, http.StatusForbidden, errors.CodeForbidden, "Forbidden")
	case errors.ErrEmailSendFailed:
		return SendError(c, http.StatusInternalServerError, errors.CodeInternalServer, "Failed to send email")
	case errors.ErrDatabaseConnection:
		return SendError(c, http.StatusInternalServerError, errors.CodeInternalServer, "Database connection error")
	default:
		return SendError(c, http.StatusInternalServerError, errors.CodeInternalServer, "Internal server error", err.Error())
	}
}

// getHTTPStatusFromErrorCode maps error codes to HTTP status codes
func getHTTPStatusFromErrorCode(code string) int {
	switch code {
	case errors.CodeValidation, errors.CodeBadRequest:
		return http.StatusBadRequest
	case errors.CodeUnauthorized:
		return http.StatusUnauthorized
	case errors.CodeForbidden:
		return http.StatusForbidden
	case errors.CodeNotFound:
		return http.StatusNotFound
	case errors.CodeConflict:
		return http.StatusConflict
	case errors.CodeInternalServer:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
} 
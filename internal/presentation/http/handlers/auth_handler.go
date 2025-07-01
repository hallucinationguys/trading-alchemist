package handlers

import (
	"github.com/gofiber/fiber/v2"

	"trading-alchemist/internal/application/dto"
	"trading-alchemist/internal/application/usecases"
	"trading-alchemist/internal/presentation/responses"
	"trading-alchemist/pkg/errors"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authUseCase *usecases.AuthUseCase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUseCase *usecases.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// SendMagicLink sends a magic link to user's email
// @Summary Send magic link
// @Description Sends a magic link to the specified email address for passwordless authentication. The magic link will be valid for the configured TTL period (default 15 minutes).
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.SendMagicLinkRequest true "Email address to send magic link to"
// @Success 200 {object} responses.SuccessResponse{data=dto.SendMagicLinkResponse} "Magic link sent successfully"
// @Failure 400 {object} responses.ErrorResponse "Invalid email address or validation error"
// @Failure 429 {object} responses.ErrorResponse "Too many requests - rate limit exceeded"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /auth/magic-link [post]
func (h *AuthHandler) SendMagicLink(c *fiber.Ctx) error {
	var req dto.SendMagicLinkRequest
	
	if err := c.BodyParser(&req); err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	// Set IP address and user agent from request
	req.IPAddress = c.IP()
	req.UserAgent = c.Get("User-Agent")

	// Call use case
	response, err := h.authUseCase.SendMagicLink(c.Context(), &req)
	if err != nil {
		// Handle different error types
		if appErr, ok := err.(*errors.AppError); ok {
			switch appErr.Code {
			case errors.CodeValidation:
				return responses.SendError(c, fiber.StatusBadRequest, string(appErr.Code), appErr.Message)
			case errors.CodeNotFound:
				return responses.SendError(c, fiber.StatusNotFound, string(appErr.Code), appErr.Message)
			default:
				return responses.SendError(c, fiber.StatusInternalServerError, string(appErr.Code), appErr.Message)
			}
		}
		return responses.SendError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred")
	}

	return responses.SendSuccess(c, response, "Magic link sent successfully")
}

// VerifyMagicLink verifies a magic link token and returns JWT tokens
// @Summary Verify magic link
// @Description Verifies a magic link token from an email link and returns a JWT access token if valid. The magic link token is consumed and cannot be used again.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.VerifyMagicLinkRequest true "Magic link token to verify"
// @Success 200 {object} responses.SuccessResponse{data=dto.VerifyMagicLinkResponse} "Magic link verified successfully"
// @Failure 400 {object} responses.ErrorResponse "Invalid token format or token missing"
// @Failure 401 {object} responses.ErrorResponse "Invalid, expired, or already used token"
// @Failure 404 {object} responses.ErrorResponse "Token not found"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /auth/verify [post]
func (h *AuthHandler) VerifyMagicLink(c *fiber.Ctx) error {
	var req dto.VerifyMagicLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	if req.Token == "" {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Token is required")
	}

	// Call use case
	response, err := h.authUseCase.VerifyMagicLink(c.Context(), &req)
	if err != nil {
		// Handle different error types
		if appErr, ok := err.(*errors.AppError); ok {
			switch appErr.Code {
			case errors.CodeValidation:
				return responses.SendError(c, fiber.StatusBadRequest, string(appErr.Code), appErr.Message)
			case errors.CodeNotFound:
				return responses.SendError(c, fiber.StatusNotFound, string(appErr.Code), appErr.Message)
			case errors.CodeUnauthorized:
				return responses.SendError(c, fiber.StatusUnauthorized, string(appErr.Code), appErr.Message)
			case errors.CodeForbidden:
				return responses.SendError(c, fiber.StatusForbidden, string(appErr.Code), appErr.Message)
			default:
				return responses.SendError(c, fiber.StatusInternalServerError, string(appErr.Code), appErr.Message)
			}
		}
		return responses.SendError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred")
	}

	return responses.SendSuccess(c, response, "Authentication successful")
} 
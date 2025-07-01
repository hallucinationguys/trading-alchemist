package handlers

import (
	"github.com/gofiber/fiber/v2"

	"trading-alchemist/internal/application/dto"
	"trading-alchemist/internal/presentation/responses"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	// In a real implementation, this would have dependencies like:
	// authUseCase application.AuthUseCase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
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

	// TODO: Add validation
	// TODO: Implement actual magic link sending logic
	
	response := dto.SendMagicLinkResponse{
		Message: "If this email is registered, a magic link has been sent",
		Sent:    true,
	}

	return responses.SendSuccess(c, response, "Magic link sent successfully")
}

// VerifyMagicLink verifies a magic link token and returns JWT tokens
// @Summary Verify magic link
// @Description Verifies a magic link token and returns JWT access token if valid. The magic link token is consumed and cannot be used again.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.VerifyMagicLinkRequest true "Magic link token to verify"
// @Success 200 {object} responses.SuccessResponse{data=dto.VerifyMagicLinkResponse} "Magic link verified successfully"
// @Failure 400 {object} responses.ErrorResponse "Invalid token format"
// @Failure 401 {object} responses.ErrorResponse "Invalid, expired, or already used token"
// @Failure 404 {object} responses.ErrorResponse "Token not found"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /auth/verify [post]
func (h *AuthHandler) VerifyMagicLink(c *fiber.Ctx) error {
	var req dto.VerifyMagicLinkRequest
	
	if err := c.BodyParser(&req); err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	// TODO: Add validation
	// TODO: Implement actual token verification logic
	
	response := dto.VerifyMagicLinkResponse{
		User: dto.UserResponse{
			// Mock user data - replace with actual user from database
			Email:         "user@example.com",
			EmailVerified: true,
			FullName:      "John Doe",
			DisplayName:   "John Doe",
			IsActive:      true,
		},
		AccessToken: "mock-jwt-token",
		TokenType:   "Bearer",
		ExpiresIn:   86400, // 24 hours in seconds
	}

	return responses.SendSuccess(c, response, "Authentication successful")
} 
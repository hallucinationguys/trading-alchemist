package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"trading-alchemist/internal/application/auth"
	"trading-alchemist/internal/presentation/responses"
	"trading-alchemist/pkg/utils"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userUseCase *auth.UserUseCase
	authUseCase *auth.AuthUseCase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase *auth.UserUseCase, authUseCase *auth.AuthUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		authUseCase: authUseCase,
	}
}

// GetProfile retrieves the current user's profile
// @Summary Get user profile
// @Description Retrieves the profile information for the currently authenticated user. Requires valid JWT token in Authorization header.
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} responses.SuccessResponse{data=auth.GetUserResponse} "User profile retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized - invalid or missing token"
// @Failure 404 {object} responses.ErrorResponse "User not found"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	// Extract user from context (set by auth middleware)
	userClaims, ok := c.Locals("user").(*utils.Claims)
	if !ok || userClaims == nil {
		return responses.SendError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
	}

	userID, err := uuid.Parse(userClaims.Subject)
	if err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid user ID format")
	}

	// Call use case to get profile
	userProfile, err := h.userUseCase.GetUserProfile(c.Context(), userID)
	if err != nil {
		return responses.HandleError(c, err)
	}

	response := auth.GetUserResponse{
		User: *userProfile,
	}

	return responses.SendSuccess(c, response, "Profile retrieved successfully")
}

// UpdateProfile updates the current user's profile
// @Summary Update user profile
// @Description Updates the profile information for the currently authenticated user. Only provided fields will be updated (partial update).
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body auth.UpdateUserRequest true "User profile information to update"
// @Success 200 {object} responses.SuccessResponse{data=auth.UpdateUserResponse} "User profile updated successfully"
// @Failure 400 {object} responses.ErrorResponse "Invalid request data or validation error"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized - invalid or missing token"
// @Failure 404 {object} responses.ErrorResponse "User not found"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	var req auth.UpdateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	// Validate the request fields
	if req.FirstName != nil && !utils.ValidateStringLength(*req.FirstName, 1, 100) {
		return responses.SendError(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "First name must be between 1 and 100 characters")
	}
	if req.LastName != nil && !utils.ValidateStringLength(*req.LastName, 1, 100) {
		return responses.SendError(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Last name must be between 1 and 100 characters")
	}

	// Extract user from context
	userClaims, ok := c.Locals("user").(*utils.Claims)
	if !ok || userClaims == nil {
		return responses.SendError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
	}

	userID, err := uuid.Parse(userClaims.Subject)
	if err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid user ID format")
	}

	// Call use case to update profile
	updatedUser, err := h.userUseCase.UpdateUserProfile(c.Context(), userID, &req)
	if err != nil {
		return responses.HandleError(c, err)
	}

	response := auth.UpdateUserResponse{
		User: *updatedUser,
	}

	return responses.SendSuccess(c, response, "Profile updated successfully")
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func getFullName(firstName, lastName *string) string {
	if firstName != nil && lastName != nil {
		return *firstName + " " + *lastName
	}
	if firstName != nil {
		return *firstName
	}
	if lastName != nil {
		return *lastName
	}
	return ""
}

func getDisplayName(firstName, lastName *string, email string) string {
	fullName := getFullName(firstName, lastName)
	if fullName != "" {
		return fullName
	}
	return email
} 
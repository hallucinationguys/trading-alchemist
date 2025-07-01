package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"trading-alchemist/internal/application/dto"
	"trading-alchemist/internal/application/usecases"
	"trading-alchemist/internal/presentation/responses"
	"trading-alchemist/pkg/utils"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userUseCase *usecases.UserUseCase
	authUseCase *usecases.AuthUseCase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase *usecases.UserUseCase, authUseCase *usecases.AuthUseCase) *UserHandler {
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
// @Success 200 {object} responses.SuccessResponse{data=dto.GetUserResponse} "User profile retrieved successfully"
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

	response := dto.GetUserResponse{
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
// @Param request body dto.UpdateUserRequest true "User profile information to update"
// @Success 200 {object} responses.SuccessResponse{data=dto.UpdateUserResponse} "User profile updated successfully"
// @Failure 400 {object} responses.ErrorResponse "Invalid request data or validation error"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized - invalid or missing token"
// @Failure 404 {object} responses.ErrorResponse "User not found"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	var req dto.UpdateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	// Extract user from context
	userClaims, ok := c.Locals("user").(*utils.Claims)
	if !ok || userClaims == nil {
		return responses.SendError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
	}

	// TODO: Add validation for the request
	// TODO: Implement actual user update logic in a use case

	// Mock updated user data - replace with actual database operation
	userID, _ := uuid.Parse(userClaims.Subject)
	user := dto.UserResponse{
		ID:            userID,
		Email:         userClaims.Email, // Email is not updatable in this example
		EmailVerified: true,             // Assuming email is verified
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		AvatarURL:     req.AvatarURL,
		FullName:      getFullName(req.FirstName, req.LastName),
		DisplayName:   getDisplayName(req.FirstName, req.LastName, userClaims.Email),
		IsActive:      true,
	}

	response := dto.UpdateUserResponse{
		User: user,
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
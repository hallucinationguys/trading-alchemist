package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"trading-alchemist/internal/application/dto"
	"trading-alchemist/internal/presentation/responses"
)

// UserHandler handles user-related requests
type UserHandler struct {
	// In a real implementation, this would have dependencies like:
	// userUseCase application.UserUseCase
}

// NewUserHandler creates a new user handler
func NewUserHandler() *UserHandler {
	return &UserHandler{}
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
	// TODO: Extract user ID from JWT token
	// TODO: Implement actual user retrieval logic
	
	// Mock user data - replace with actual user from database
	user := dto.UserResponse{
		ID:            uuid.New(),
		Email:         "user@example.com",
		EmailVerified: true,
		FirstName:     stringPtr("John"),
		LastName:      stringPtr("Doe"),
		FullName:      "John Doe",
		DisplayName:   "John Doe",
		IsActive:      true,
	}

	response := dto.GetUserResponse{
		User: user,
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

	// TODO: Extract user ID from JWT token
	// TODO: Add validation
	// TODO: Implement actual user update logic
	
	// Mock updated user data - replace with actual database operation
	user := dto.UserResponse{
		ID:            uuid.New(),
		Email:         "user@example.com",
		EmailVerified: true,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		AvatarURL:     req.AvatarURL,
		FullName:      getFullName(req.FirstName, req.LastName),
		DisplayName:   getDisplayName(req.FirstName, req.LastName, "user@example.com"),
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
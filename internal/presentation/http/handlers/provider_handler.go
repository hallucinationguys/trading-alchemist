package handlers

import (
	"log"
	"trading-alchemist/internal/application/chat"
	"trading-alchemist/internal/presentation/responses"
	"trading-alchemist/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ProviderHandler handles provider-related and user-provider-setting-related requests.
type ProviderHandler struct {
	providerUseCase           *chat.UserProviderSettingUseCase
	modelAvailabilityUseCase  *chat.ModelAvailabilityUseCase
}

// NewProviderHandler creates a new ProviderHandler.
func NewProviderHandler(
	providerUseCase *chat.UserProviderSettingUseCase,
	modelAvailabilityUseCase *chat.ModelAvailabilityUseCase,
) *ProviderHandler {
	return &ProviderHandler{
		providerUseCase:          providerUseCase,
		modelAvailabilityUseCase: modelAvailabilityUseCase,
	}
}

// ListProviders retrieves a list of all available providers.
// @Summary List available providers
// @Description Retrieves a list of all active LLM providers supported by the system.
// @Tags Providers
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} responses.SuccessResponse{data=[]chat.ProviderResponse} "Providers retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /providers [get]
func (h *ProviderHandler) ListProviders(c *fiber.Ctx) error {
	providers, err := h.providerUseCase.ListProviders(c.Context())
	if err != nil {
		return responses.HandleError(c, err)
	}
	return responses.SendSuccess(c, providers, "Providers retrieved successfully")
}

// ListUserSettings retrieves the provider settings for the current user.
// @Summary List user's provider settings
// @Description Retrieves all provider settings for the currently authenticated user.
// @Tags Providers
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} responses.SuccessResponse{data=[]chat.UserProviderSettingResponse} "User provider settings retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /providers/settings [get]
func (h *ProviderHandler) ListUserSettings(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(*utils.Claims)
	userID, err := uuid.Parse(userClaims.Subject)
	if err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid user ID format in token")
	}

	settings, err := h.providerUseCase.ListUserSettings(c.Context(), userID)
	if err != nil {
		return responses.HandleError(c, err)
	}
	return responses.SendSuccess(c, settings, "User provider settings retrieved successfully")
}

// UpsertUserSetting creates or updates a provider setting for the current user.
// @Summary Create or update a provider setting
// @Description Creates or updates a provider setting (API key, base URL) for the authenticated user.
// @Tags Providers
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body chat.UpsertUserProviderSettingRequest true "Provider setting information"
// @Success 200 {object} responses.SuccessResponse{data=chat.UserProviderSettingResponse} "Provider setting saved successfully"
// @Failure 400 {object} responses.ErrorResponse "Invalid request body or validation error"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 404 {object} responses.ErrorResponse "Provider not found"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /providers/settings [post]
func (h *ProviderHandler) UpsertUserSetting(c *fiber.Ctx) error {
	var req chat.UpsertUserProviderSettingRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	log.Printf("Received provider setting request: ProviderID=%s, APIKey length=%d", req.ProviderID, len(req.APIKey))

	// Validate required fields
	if req.ProviderID == uuid.Nil {
		return responses.SendError(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Provider ID is required")
	}

	userClaims := c.Locals("user").(*utils.Claims)
	userID, err := uuid.Parse(userClaims.Subject)
	if err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid user ID format in token")
	}

	log.Printf("Processing provider setting for user: %s", userID)

	// Note: API key validation is now handled in the use case
	// - For new settings: API key is required
	// - For updates: API key is optional (empty means keep existing)
	setting, err := h.providerUseCase.UpsertUserProviderSetting(c.Context(), userID, &req)
	if err != nil {
		log.Printf("Failed to upsert provider setting: %v", err)
		return responses.HandleError(c, err)
	}

	log.Printf("Successfully saved provider setting for user: %s", userID)
	return responses.SendSuccess(c, setting, "Provider setting saved successfully")
}

// GetAvailableModels retrieves available models with API key status for the user
// @Summary Get available models with API key status
// @Description Retrieves all available models with their API key configuration status in a single optimized call
// @Tags Providers
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} responses.SuccessResponse{data=[]chat.ProviderResponse} "Available models retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /providers/available-models [get]
func (h *ProviderHandler) GetAvailableModels(c *fiber.Ctx) error {
	userClaims := c.Locals("user").(*utils.Claims)
	userID, err := uuid.Parse(userClaims.Subject)
	if err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid user ID format in token")
	}

	models, err := h.modelAvailabilityUseCase.GetAvailableModelsWithAPIKeyStatus(c.Context(), userID)
	if err != nil {
		return responses.HandleError(c, err)
	}
	
	return responses.SendSuccess(c, models, "Available models retrieved successfully")
} 
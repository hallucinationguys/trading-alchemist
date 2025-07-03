package dto

import (
	"time"

	"trading-alchemist/internal/domain/entities"

	"github.com/google/uuid"
)

// --- Response DTOs ---

// ProviderResponse represents a publicly available provider that the system supports.
type ProviderResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	IsActive    bool      `json:"is_active"`
	Models      []ModelResponse `json:"models"`
}

// UserProviderSettingResponse represents a user's configured provider setting.
// It combines information about the provider and the user's specific configuration.
type UserProviderSettingResponse struct {
	ID                  uuid.UUID `json:"id"`
	ProviderID          uuid.UUID `json:"provider_id"`
	ProviderName        string    `json:"provider_name"`
	ProviderDisplayName string    `json:"provider_display_name"`
	APIKeySet           bool      `json:"api_key_set"` // Indicates if the API key is configured, without exposing the key.
	APIBaseOverride     *string   `json:"api_base_override,omitempty"`
	IsActive            bool      `json:"is_active"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// ModelResponse represents a single model for a provider.
type ModelResponse struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	DisplayName string  `json:"display_name"`
	IsActive   bool     `json:"is_active"`
	Tags       []string `json:"tags"` // Example: ["LLM", "CHAT", "200K"]
}

// --- Request DTOs ---

// UpsertUserProviderSettingRequest is used to create or update a user's provider setting.
type UpsertUserProviderSettingRequest struct {
	ProviderID      uuid.UUID `json:"provider_id" validate:"required"`
	APIKey          string    `json:"api_key" validate:"required"`
	APIBaseOverride *string   `json:"api_base_override,omitempty" validate:"omitempty,url"`
	IsActive        *bool     `json:"is_active,omitempty"`
}

// --- Helper Functions ---

// ToProviderResponse converts a provider entity to a response DTO.
func ToProviderResponse(p *entities.Provider) ProviderResponse {
	return ProviderResponse{
		ID:          p.ID,
		Name:        p.Name,
		DisplayName: p.DisplayName,
		IsActive:    p.IsActive,
		Models:      ToModelResponses(p.Models),
	}
}

// ToUserProviderSettingResponse converts a setting entity and its associated provider entity to a response DTO.
func ToUserProviderSettingResponse(setting *entities.UserProviderSetting, provider *entities.Provider) UserProviderSettingResponse {
	return UserProviderSettingResponse{
		ID:                  setting.ID,
		ProviderID:          setting.ProviderID,
		ProviderName:        provider.Name,
		ProviderDisplayName: provider.DisplayName,
		APIKeySet:           setting.EncryptedAPIKey != nil && *setting.EncryptedAPIKey != "",
		APIBaseOverride:     setting.APIBaseOverride,
		IsActive:            setting.IsActive,
		UpdatedAt:           setting.UpdatedAt,
	}
}

// ToModelResponse converts an entities.Model to a ModelResponse DTO.
func ToModelResponse(m *entities.Model) ModelResponse {
	// Logic to create tags can be more sophisticated based on model properties
	tags := []string{"LLM", "CHAT"}
	if m.SupportsVision {
		tags = append(tags, "VISION")
	}

	return ModelResponse{
		ID:          m.ID.String(),
		Name:        m.Name,
		DisplayName: m.DisplayName,
		IsActive:    m.IsActive,
		Tags:        tags,
	}
}

// ToModelResponses converts a slice of entities.Model to a slice of ModelResponse DTOs.
func ToModelResponses(models []*entities.Model) []ModelResponse {
	if models == nil {
		return []ModelResponse{}
	}
	responses := make([]ModelResponse, len(models))
	for i, m := range models {
		responses[i] = ToModelResponse(m)
	}
	return responses
} 
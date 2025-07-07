package usecases

import (
	"context"
	"fmt"
	"log"
	"strings"
	"trading-alchemist/internal/application/dto"
	"trading-alchemist/internal/config"
	"trading-alchemist/internal/domain/entities"
	"trading-alchemist/internal/domain/services"
	"trading-alchemist/internal/infrastructure/database"
	"trading-alchemist/internal/infrastructure/llm/prompts"
	"trading-alchemist/pkg/errors"
	"trading-alchemist/pkg/utils"

	"github.com/google/uuid"
)

// ConversationUseCase handles the business logic for conversation management.
type ConversationUseCase struct {
	dbService     *database.Service
	config        *config.Config
	llmService    services.LLMService
	promptManager *prompts.PromptManager
}

// NewConversationUseCase creates a new ConversationUseCase instance.
func NewConversationUseCase(
	dbService *database.Service,
	config *config.Config,
	llmService services.LLMService,
) *ConversationUseCase {
	return &ConversationUseCase{
		dbService:     dbService,
		config:        config,
		llmService:    llmService,
		promptManager: prompts.NewPromptManager(),
	}
}

// CreateConversation creates a new chat conversation.
func (uc *ConversationUseCase) CreateConversation(ctx context.Context, req *dto.CreateConversationRequest) (*dto.ConversationDetailResponse, error) {
	var createdConv *entities.Conversation

	err := uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		var targetModel *entities.Model
		var err error

		if req.ModelName != nil && *req.ModelName != "" {
			// Use the specific model requested
			parts := strings.Split(*req.ModelName, "/")
			if len(parts) != 2 {
				return errors.NewAppError(errors.CodeValidation, fmt.Sprintf("Invalid model format: %s. Expected 'provider/model_name'", *req.ModelName), nil)
			}
			providerName, modelName := parts[0], parts[1]

			// Find the provider
			targetProvider, err := provider.Provider().GetByName(ctx, providerName)
			if err != nil {
				if err == errors.ErrProviderNotFound {
					return errors.NewAppError(errors.CodeNotFound, fmt.Sprintf("Provider '%s' not found", providerName), err)
				}
				return fmt.Errorf("failed to find provider '%s': %w", providerName, err)
			}

			// Find the model within the provider
			targetModel, err = provider.Model().GetModelByName(ctx, targetProvider.ID, modelName)
			if err != nil {
				if err == errors.ErrModelNotFound {
					return errors.NewAppError(errors.CodeNotFound, fmt.Sprintf("Model '%s' not found for provider '%s'", modelName, providerName), err)
				}
				return fmt.Errorf("failed to find model '%s': %w", modelName, err)
			}
		} else {
			// Use default model from config if no specific model requested
			if uc.config.App.DefaultModel == "" {
				return errors.NewAppError(errors.CodeConfiguration, "No model specified and no default model configured", nil)
			}

			parts := strings.Split(uc.config.App.DefaultModel, "/")
			if len(parts) != 2 {
				return errors.NewAppError(errors.CodeConfiguration, fmt.Sprintf("Invalid default model format in config: %s", uc.config.App.DefaultModel), nil)
			}
			providerName, modelName := parts[0], parts[1]

			// Find the provider
			targetProvider, err := provider.Provider().GetByName(ctx, providerName)
			if err != nil {
				return fmt.Errorf("failed to find default provider '%s': %w", providerName, err)
			}

			// Find the model within the provider
			targetModel, err = provider.Model().GetModelByName(ctx, targetProvider.ID, modelName)
			if err != nil {
				return fmt.Errorf("failed to find default model '%s': %w", modelName, err)
			}
		}

		// Create Conversation
		newConv := &entities.Conversation{
			UserID:  req.UserID,
			Title:   req.Title,
			ModelID: targetModel.ID,
		}
		createdConv, err = provider.Conversation().Create(ctx, newConv)
		if err != nil {
			return fmt.Errorf("failed to create conversation: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.ConversationDetailResponse{
		ID:      createdConv.ID,
		Title:   createdConv.Title,
		ModelID: createdConv.ModelID,
	}, nil
}

// GetUserConversations retrieves a paginated list of conversations for a specific user.
func (uc *ConversationUseCase) GetUserConversations(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*dto.ConversationSummaryResponse, error) {
	var conversations []*entities.Conversation

	err := uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		var err error
		conversations, err = provider.Conversation().GetByUserID(ctx, userID, limit, offset)
		if err != nil {
			return fmt.Errorf("failed to get user conversations: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	response := make([]*dto.ConversationSummaryResponse, len(conversations))
	for i, conv := range conversations {
		response[i] = &dto.ConversationSummaryResponse{
			ID:            conv.ID,
			Title:         conv.Title,
			LastMessageAt: conv.LastMessageAt,
			ModelID:       conv.ModelID,
		}
	}

	return response, nil
}

// GetConversationDetails retrieves the full details of a single conversation, including its messages.
func (uc *ConversationUseCase) GetConversationDetails(ctx context.Context, conversationID, userID uuid.UUID) (*dto.ConversationDetailResponse, error) {
	var conversation *entities.Conversation
	var messages []*entities.Message
	var artifactsMap map[uuid.UUID][]*entities.Artifact

	err := uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		var err error
		conversation, err = provider.Conversation().GetByID(ctx, conversationID)
		if err != nil {
			return fmt.Errorf("failed to get conversation: %w", err)
		}

		// Security check: ensure the user owns the conversation
		if conversation.UserID != userID {
			return errors.ErrForbidden
		}

		// Fetch last 100 messages for now (should be paginated in production)
		messages, err = provider.Message().GetByConversationID(ctx, conversationID, 100, 0)
		if err != nil {
			return fmt.Errorf("failed to get messages: %w", err)
		}

		// Fetch artifacts for all messages
		artifactsMap = make(map[uuid.UUID][]*entities.Artifact)
		for _, msg := range messages {
			artifacts, err := provider.Artifact().GetByMessageID(ctx, msg.ID)
			if err != nil {
				return fmt.Errorf("failed to get artifacts for message %s: %w", msg.ID, err)
			}
			if len(artifacts) > 0 {
				artifactsMap[msg.ID] = artifacts
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert messages to DTOs
	messageDTOs := make([]dto.MessageResponse, len(messages))
	for i, msg := range messages {
		messageDTOs[i] = dto.MessageResponse{
			ID:        msg.ID,
			Role:      string(msg.Role),
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
			Artifacts: toArtifactResponses(artifactsMap[msg.ID]),
		}
	}

	return &dto.ConversationDetailResponse{
		ID:           conversation.ID,
		Title:        conversation.Title,
		ModelID:      conversation.ModelID,
		SystemPrompt: conversation.SystemPrompt,
		Messages:     messageDTOs,
	}, nil
}

// UpdateConversationTitle updates the title of a conversation.
func (uc *ConversationUseCase) UpdateConversationTitle(ctx context.Context, conversationID, userID uuid.UUID, req *dto.UpdateConversationTitleRequest) error {
	return uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		// First verify the conversation exists and user owns it
		conversation, err := provider.Conversation().GetByID(ctx, conversationID)
		if err != nil {
			if err == errors.ErrConversationNotFound {
				return errors.NewAppError(errors.CodeNotFound, "Conversation not found", err)
			}
			return fmt.Errorf("failed to get conversation: %w", err)
		}

		// Security check: ensure the user owns the conversation
		if conversation.UserID != userID {
			return errors.ErrForbidden
		}

		// Update the title
		err = provider.Conversation().UpdateTitle(ctx, conversationID, req.Title)
		if err != nil {
			return fmt.Errorf("failed to update conversation title: %w", err)
		}

		return nil
	})
}

// ArchiveConversation archives (soft deletes) a conversation.
func (uc *ConversationUseCase) ArchiveConversation(ctx context.Context, conversationID, userID uuid.UUID) error {
	return uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		// First verify the conversation exists and user owns it
		conversation, err := provider.Conversation().GetByID(ctx, conversationID)
		if err != nil {
			if err == errors.ErrConversationNotFound {
				return errors.NewAppError(errors.CodeNotFound, "Conversation not found", err)
			}
			return fmt.Errorf("failed to get conversation: %w", err)
		}

		// Security check: ensure the user owns the conversation
		if conversation.UserID != userID {
			return errors.ErrForbidden
		}

		// Archive the conversation
		err = provider.Conversation().Archive(ctx, conversationID)
		if err != nil {
			return fmt.Errorf("failed to archive conversation: %w", err)
		}

		return nil
	})
}

// GenerateConversationTitle generates a descriptive title for a conversation based on the first exchange.
// This method is called asynchronously after the first assistant response.
func (uc *ConversationUseCase) GenerateConversationTitle(ctx context.Context, conversationID uuid.UUID, userMessage, assistantMessage string) {
	log.Printf("GenerateConversationTitle called for conversation %s", conversationID)
	go func() {
		// Use a separate context for the background operation
		bgCtx := context.Background()
		
		log.Printf("Starting background title generation for conversation %s", conversationID)
		
		// Only generate title if conversation still has a default title
		var currentTitle string
		err := uc.dbService.ExecuteInTx(bgCtx, func(provider database.RepositoryProvider) error {
			conversation, err := provider.Conversation().GetByID(bgCtx, conversationID)
			if err != nil {
				return err
			}
			currentTitle = conversation.Title
			return nil
		})
		
		if err != nil {
			log.Printf("Failed to check conversation title for generation: %v", err)
			return
		}
		
		// Only generate if it's still the default title
		if currentTitle != "New Conversation" {
			log.Printf("Conversation %s already has a custom title, skipping generation", conversationID)
			return
		}
		
		// Generate title using LLM
		generatedTitle, err := uc.generateTitleWithLLM(bgCtx, conversationID, userMessage, assistantMessage)
		if err != nil {
			log.Printf("Failed to generate title with LLM for conversation %s: %v", conversationID, err)
			// Fallback to truncated user message
			generatedTitle = uc.generateFallbackTitle(userMessage)
		}
		
		// Update the conversation title
		err = uc.dbService.ExecuteInTx(bgCtx, func(provider database.RepositoryProvider) error {
			return provider.Conversation().UpdateTitle(bgCtx, conversationID, generatedTitle)
		})
		
		if err != nil {
			log.Printf("Failed to update conversation title for %s: %v", conversationID, err)
		} else {
			log.Printf("Successfully generated title for conversation %s: %s", conversationID, generatedTitle)
		}
	}()
}

// CheckShouldGenerateTitle checks if we should generate a title for the conversation.
// Returns true if this is the first exchange (2 messages total).
func (uc *ConversationUseCase) CheckShouldGenerateTitle(ctx context.Context, conversationID uuid.UUID) (bool, error) {
	var messageCount int
	err := uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		var err error
		messageCount, err = provider.Message().CountByConversationID(ctx, conversationID)
		return err
	})
	
	if err != nil {
		return false, fmt.Errorf("failed to count messages: %w", err)
	}
	
	log.Printf("Message count for conversation %s: %d", conversationID, messageCount)
	
	// Generate title after first exchange (user message + assistant response = 2 messages)
	shouldGenerate := messageCount == 2
	log.Printf("Should generate title based on count: %v (count == 2)", shouldGenerate)
	return shouldGenerate, nil
}

// CheckShouldGenerateTitleWithProvider checks if we should generate a title using an existing provider.
// This is used when we're already within a transaction to avoid nested transactions.
func (uc *ConversationUseCase) CheckShouldGenerateTitleWithProvider(provider database.RepositoryProvider, ctx context.Context, conversationID uuid.UUID) (bool, error) {
	messageCount, err := provider.Message().CountByConversationID(ctx, conversationID)
	if err != nil {
		return false, fmt.Errorf("failed to count messages: %w", err)
	}
	
	log.Printf("Message count for conversation %s: %d", conversationID, messageCount)
	
	// Generate title after first exchange (user message + assistant response = 2 messages)
	shouldGenerate := messageCount == 2
	log.Printf("Should generate title based on count: %v (count == 2)", shouldGenerate)
	return shouldGenerate, nil
}

// generateTitleWithLLM uses the LLM service to generate a descriptive title using the user's API key.
func (uc *ConversationUseCase) generateTitleWithLLM(ctx context.Context, conversationID uuid.UUID, userMessage, assistantMessage string) (string, error) {
	log.Printf("Generating title with LLM for conversation %s", conversationID)

	// Get the title generation prompt
	systemPrompt, err := uc.promptManager.GetSystemPrompt("title_generation")
	if err != nil {
		log.Printf("Failed to get title generation system prompt: %v", err)
		return uc.generateFallbackTitle(userMessage), nil
	}

	// Prepare template data
	templateData := prompts.TitleGenerationData{
		UserMessage:      strings.TrimSpace(userMessage),
		AssistantMessage: strings.TrimSpace(assistantMessage),
	}

	// Render the user prompt with the conversation data
	userPrompt, err := uc.promptManager.RenderUserPrompt("title_generation", templateData)
	if err != nil {
		log.Printf("Failed to render title generation user prompt: %v", err)
		return uc.generateFallbackTitle(userMessage), nil
	}

	// Get conversation details, user API key, and model info
	var conversation *entities.Conversation
	var titleProvider *entities.Provider
	var titleModel *entities.Model
	var userSetting *entities.UserProviderSetting
	
	err = uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		// Get conversation to find user_id
		var err error
		conversation, err = provider.Conversation().GetByID(ctx, conversationID)
		if err != nil {
			return fmt.Errorf("failed to get conversation: %w", err)
		}

		// Get OpenAI provider
		providers, err := provider.Provider().GetAll(ctx)
		if err != nil {
			return fmt.Errorf("failed to get providers: %w", err)
		}
		
		for _, p := range providers {
			if p.Name == "openai" {
				titleProvider = p
				break
			}
		}
		
		if titleProvider == nil {
			return fmt.Errorf("openai provider not found")
		}

		// Get gpt-4o-mini model
		models, err := provider.Model().GetActiveModelsByProviderID(ctx, titleProvider.ID)
		if err != nil {
			return fmt.Errorf("failed to get models for provider: %w", err)
		}
		
		for _, m := range models {
			if m.Name == "gpt-4o-mini" {
				titleModel = m
				break
			}
		}
		
		if titleModel == nil {
			return fmt.Errorf("gpt-4o-mini model not found")
		}

		// Get user's API key for OpenAI
		userSetting, err = provider.UserProviderSetting().GetByUserIDAndProviderID(ctx, conversation.UserID, titleProvider.ID)
		if err != nil {
			return fmt.Errorf("failed to get user API key for OpenAI: %w", err)
		}
		
		if !userSetting.IsActive || userSetting.EncryptedAPIKey == nil || *userSetting.EncryptedAPIKey == "" {
			return fmt.Errorf("user has no active API key for OpenAI")
		}

		return nil
	})

	if err != nil {
		log.Printf("Failed to get title generation prerequisites: %v", err)
		return uc.generateFallbackTitle(userMessage), nil
	}

	// Decrypt the API key
	encryptionKey, err := uc.config.GetEncryptionKey()
	if err != nil {
		log.Printf("Failed to get encryption key: %v", err)
		return uc.generateFallbackTitle(userMessage), nil
	}
	
	decryptedAPIKey, err := utils.Decrypt(*userSetting.EncryptedAPIKey, encryptionKey)
	if err != nil {
		log.Printf("Failed to decrypt API key: %v", err)
		return uc.generateFallbackTitle(userMessage), nil
	}

	// Create messages for title generation
	messages := []*entities.Message{
		{
			Role:    entities.MessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    entities.MessageRoleUser,
			Content: userPrompt,
		},
	}

	// Get API base override if set
	apiBaseOverride := ""
	if userSetting.APIBaseOverride != nil {
		apiBaseOverride = *userSetting.APIBaseOverride
	}

	// Call LLM service for title generation
	log.Printf("Making LLM call for title generation using user's API key")
	llmEventCh, err := uc.llmService.StreamChatCompletion(ctx, titleProvider, titleModel, messages, decryptedAPIKey, apiBaseOverride)
	if err != nil {
		log.Printf("Failed to start LLM stream for title generation: %v", err)
		return uc.generateFallbackTitle(userMessage), nil
	}

	// Collect the streaming response
	var titleContent strings.Builder
	for event := range llmEventCh {
		if event.Error != nil {
			log.Printf("Error during title generation LLM stream: %v", event.Error)
			return uc.generateFallbackTitle(userMessage), nil
		}
		
		titleContent.WriteString(event.ContentDelta)
		
		if event.IsLast {
			break
		}
	}

	generatedTitle := strings.TrimSpace(titleContent.String())
	
	// Validate and clean the generated title
	if generatedTitle == "" {
		log.Printf("LLM generated empty title, using fallback")
		return uc.generateFallbackTitle(userMessage), nil
	}

	// Remove any quotes or extra formatting
	generatedTitle = strings.Trim(generatedTitle, `"'`)
	
	// Ensure it's not too long
	if len(generatedTitle) > 50 {
		generatedTitle = generatedTitle[:47] + "..."
	}

	log.Printf("Successfully generated title: %s", generatedTitle)
	return generatedTitle, nil
}

// generateFallbackTitle creates a fallback title from the user message.
func (uc *ConversationUseCase) generateFallbackTitle(userMessage string) string {
	// Clean up the message
	title := strings.TrimSpace(userMessage)
	
	// Remove newlines and extra spaces
	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.ReplaceAll(title, "\t", " ")
	
	// Collapse multiple spaces
	for strings.Contains(title, "  ") {
		title = strings.ReplaceAll(title, "  ", " ")
	}
	
	// Truncate to 50 characters
	if len(title) > 50 {
		title = title[:47] + "..."
	}
	
	// Ensure it's not empty
	if title == "" {
		title = "New Conversation"
	}
	
	return title
}

// Helper function to convert artifacts to responses (shared with chat_usecase)
func toArtifactResponses(artifacts []*entities.Artifact) []dto.ArtifactResponse {
	if artifacts == nil {
		return nil
	}
	responses := make([]dto.ArtifactResponse, len(artifacts))
	for i, a := range artifacts {
		responses[i] = dto.ArtifactResponse{
			ID:       a.ID,
			Title:    a.Title,
			Type:     string(a.Type),
			Language: a.Language,
			Content:  a.Content,
		}
	}
	return responses
} 
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
	"trading-alchemist/pkg/errors"
	"trading-alchemist/pkg/utils"

	"github.com/google/uuid"
)

// ChatUseCase handles the business logic for chat operations.
type ChatUseCase struct {
	dbService    *database.Service
	config       *config.Config
	llmService   services.LLMService
	// Add other repositories as needed
}

// NewChatUseCase creates a new ChatUseCase.
func NewChatUseCase(
	dbService *database.Service,
	config *config.Config,
	llmService services.LLMService,
) *ChatUseCase {
	return &ChatUseCase{
		dbService:    dbService,
		config:       config,
		llmService:   llmService,
	}
}

// CreateConversation creates a new chat conversation.
func (uc *ChatUseCase) CreateConversation(ctx context.Context, req *dto.CreateConversationRequest) (*dto.ConversationDetailResponse, error) {
	var createdConv *entities.Conversation

	err := uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		// For now, let's use a simpler approach - use a default model from config
		// This will be enhanced later to handle user's configured models
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

		// 2. Create Conversation
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
func (uc *ChatUseCase) GetUserConversations(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*dto.ConversationSummaryResponse, error) {
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
func (uc *ChatUseCase) GetConversationDetails(ctx context.Context, conversationID, userID uuid.UUID) (*dto.ConversationDetailResponse, error) {
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

		// For simplicity, we'll fetch the last 100 messages.
		// In a real app, this should be paginated.
		messages, err = provider.Message().GetByConversationID(ctx, conversationID, 100, 0)
		if err != nil {
			return fmt.Errorf("failed to get messages: %w", err)
		}

		// Fetch artifacts for all messages in one go
		artifactsMap = make(map[uuid.UUID][]*entities.Artifact)
		for _, msg := range messages {
			artifacts, err := provider.Artifact().GetByMessageID(ctx, msg.ID)
			if err != nil {
				// We can decide to either fail the whole request or just log the error
				// For now, let's fail, but in a real app, logging might be better.
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

// PostMessage adds a new message to a conversation and starts a streaming LLM response.
// It returns a channel that the handler can use to stream events to the client.
func (uc *ChatUseCase) PostMessage(ctx context.Context, conversationID, userID uuid.UUID, req *dto.PostMessageRequest) (<-chan services.ChatStreamEvent, error) {
	var conversationHistory []*entities.Message
	var convProvider *entities.Provider
	var convModel *entities.Model
	var userSetting *entities.UserProviderSetting

	err := uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		// 1. Get conversation and verify ownership
		conversation, err := provider.Conversation().GetByID(ctx, conversationID)
		if err != nil {
			return fmt.Errorf("failed to get conversation: %w", err)
		}
		if conversation.UserID != userID {
			return errors.ErrForbidden
		}

		// 1a. Get the model and provider for this message
		// If a specific model is requested in the message, use that; otherwise use conversation's default model
		modelID := conversation.ModelID
		if req.ModelID != nil {
			modelID = *req.ModelID
		}
		
		convModel, err = provider.Model().GetByID(ctx, modelID)
		if err != nil {
			return fmt.Errorf("failed to get model: %w", err)
		}
		convProvider, err = provider.Provider().GetByID(ctx, convModel.ProviderID)
		if err != nil {
			return fmt.Errorf("failed to get provider for model: %w", err)
		}

		// 1b. Get User Provider Settings
		userSetting, err = provider.UserProviderSetting().GetByUserIDAndProviderID(ctx, userID, convProvider.ID)
		if err != nil {
			if err == errors.ErrUserProviderSettingNotFound {
				return errors.NewAppError(errors.CodeConfiguration, fmt.Sprintf("API key for provider '%s' is not configured. Please add it in settings.", convProvider.DisplayName), err)
			}
			return fmt.Errorf("failed to get user provider settings: %w", err)
		}
		if !userSetting.IsActive || userSetting.EncryptedAPIKey == nil || *userSetting.EncryptedAPIKey == "" {
			return errors.NewAppError(errors.CodeConfiguration, fmt.Sprintf("API key for provider '%s' is not active or not set.", convProvider.DisplayName), nil)
		}

		// 2. Create the new user message
		newMessage := &entities.Message{
			ConversationID: conversationID,
			Role:           entities.MessageRoleUser,
			Content:        req.Content,
			ModelID:        &modelID, // Store which model was used for this message
		}
		createdMessage, err := provider.Message().Create(ctx, newMessage)
		if err != nil {
			return fmt.Errorf("failed to create message: %w", err)
		}

		// 3. Create any associated artifacts for the user message
		// Note: The response to the client won't include these in the initial POST response,
		// but they will be part of the conversation history for future gets.
		for _, artifactReq := range req.Artifacts {
			newArtifact := &entities.Artifact{
				MessageID: createdMessage.ID,
				Title:     artifactReq.Title,
				Type:      entities.ArtifactType(artifactReq.Type),
				Language:  artifactReq.Language,
				Content:   artifactReq.Content,
			}
			if _, err := provider.Artifact().Create(ctx, newArtifact); err != nil {
				return fmt.Errorf("failed to create artifact: %w", err)
			}
		}

		// 4. Update the conversation's last_message_at timestamp
		err = provider.Conversation().UpdateLastMessageAt(ctx, conversationID, createdMessage.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to update conversation timestamp: %w", err)
		}

		// 5. Get conversation history for LLM
		// Fetch last 20 messages for context, should be configurable
		conversationHistory, err = provider.Message().GetByConversationID(ctx, conversationID, 20, 0)
		if err != nil {
			return fmt.Errorf("failed to get conversation history: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// This channel will be returned to the handler for streaming to the client.
	clientEventChannel := make(chan services.ChatStreamEvent)

	// Decrypt API Key before starting goroutine
	encryptionKey, err := uc.config.GetEncryptionKey()
	if err != nil {
		// Create a channel to send a single error event and then close it.
		errorChan := make(chan services.ChatStreamEvent, 1)
		errorChan <- services.ChatStreamEvent{Error: fmt.Errorf("failed to get encryption key: %w", err), IsLast: true}
		close(errorChan)
		return errorChan, nil
	}
	
	decryptedAPIKey, err := utils.Decrypt(*userSetting.EncryptedAPIKey, encryptionKey)
	if err != nil {
		// Create a channel to send a single error event and then close it.
		errorChan := make(chan services.ChatStreamEvent, 1)
		errorChan <- services.ChatStreamEvent{Error: fmt.Errorf("failed to decrypt API key: %w", err), IsLast: true}
		close(errorChan)
		return errorChan, nil
	}

	apiBaseOverride := ""
	if userSetting.APIBaseOverride != nil {
		apiBaseOverride = *userSetting.APIBaseOverride
	}

	// This part happens outside the transaction
	// 6. Start LLM stream and process response in a separate goroutine
	go uc.processLLMStream(context.Background(), convProvider, convModel, conversationID, conversationHistory, clientEventChannel, decryptedAPIKey, apiBaseOverride)

	return clientEventChannel, nil
}

func (uc *ChatUseCase) processLLMStream(ctx context.Context, llmProvider *entities.Provider, llmModel *entities.Model, conversationID uuid.UUID, messages []*entities.Message, clientEventChannel chan<- services.ChatStreamEvent, apiKey, apiBaseOverride string) {
	defer close(clientEventChannel)

	llmEventCh, err := uc.llmService.StreamChatCompletion(ctx, llmProvider, llmModel, messages, apiKey, apiBaseOverride)
	if err != nil {
		log.Printf("Error starting LLM stream for conversation %s: %v", conversationID, err)
		clientEventChannel <- services.ChatStreamEvent{Error: err, IsLast: true}
		return
	}

	var responseContent strings.Builder

	for event := range llmEventCh {
		// Forward the event to the client-facing channel
		clientEventChannel <- event

		if event.Error != nil {
			log.Printf("Error during LLM stream for conversation %s: %v", conversationID, event.Error)
			// The error has been forwarded, so we just stop processing.
			return
		}
		
		responseContent.WriteString(event.ContentDelta)

		if event.IsLast {
			break // Exit the loop cleanly after the last event
		}
	}
	
	log.Printf("LLM stream finished for conversation %s. Full response: %s", conversationID, responseContent.String())
	
	// TODO: Handle tool calls from a final event if needed

	// Save the assistant's message
	assistantMessage := &entities.Message{
		ID:             uuid.New(), // Generate ID upfront to potentially send it with the stream
		ConversationID: conversationID,
		Role:           entities.MessageRoleAssistant,
		Content:        responseContent.String(),
	}

	err = uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		createdMsg, err := provider.Message().Create(ctx, assistantMessage)
		if err != nil {
			return fmt.Errorf("failed to save assistant message: %w", err)
		}
		err = provider.Conversation().UpdateLastMessageAt(ctx, conversationID, createdMsg.CreatedAt)
		if err != nil {
			// Log this error but don't fail the whole operation, as the message is already saved.
			log.Printf("Failed to update conversation timestamp after assistant message for conversation %s: %v", conversationID, err)
		}
		return nil
	})

	if err != nil {
		log.Printf("Failed to save assistant's response for conversation %s: %v", conversationID, err)
	}
}

// GetAvailableTools retrieves all available tools, optionally filtered by a provider.
func (uc *ChatUseCase) GetAvailableTools(ctx context.Context, providerID *uuid.UUID) ([]dto.ToolResponse, error) {
	var tools []*entities.Tool
	err := uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		var err error
		tools, err = provider.Tool().GetAvailableTools(ctx, providerID)
		if err != nil {
			return fmt.Errorf("failed to get available tools: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	response := make([]dto.ToolResponse, len(tools))
	for i, t := range tools {
		response[i] = dto.ToolResponse{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			Schema:      (dto.JSONB)(t.Schema),
		}
	}

	return response, nil
}

// TODO: Implement methods like:
// - GetUserConversations(ctx context.Context, userID uuid.UUID) ([]*dto.ConversationSummaryResponse, error)
// - GetConversationDetails(ctx context.Context, conversationID uuid.UUID) (*dto.ConversationDetailResponse, error)
// - PostMessage(ctx context.Context, conversationID uuid.UUID, req *dto.PostMessageRequest) (*dto.MessageResponse, error) 
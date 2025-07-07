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

// Helper function for min operation
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ChatUseCase handles the business logic for chat operations.
type ChatUseCase struct {
	dbService           *database.Service
	config              *config.Config
	llmService          services.LLMService
	conversationUseCase *ConversationUseCase
}

// NewChatUseCase creates a new ChatUseCase instance.
func NewChatUseCase(
	dbService *database.Service,
	config *config.Config,
	llmService services.LLMService,
	conversationUseCase *ConversationUseCase,
) *ChatUseCase {
	return &ChatUseCase{
		dbService:           dbService,
		config:              config,
		llmService:          llmService,
		conversationUseCase: conversationUseCase,
	}
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
	
	// Save the assistant's message
	assistantMessage := &entities.Message{
		ID:             uuid.New(), // Generate ID upfront to potentially send it with the stream
		ConversationID: conversationID,
		Role:           entities.MessageRoleAssistant,
		Content:        responseContent.String(),
	}

	var shouldGenerateTitle bool
	var userMessage string
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
		
		// Check if we should generate a title (first exchange complete)
		log.Printf("Checking if should generate title for conversation %s", conversationID)
		shouldGenerate, err := uc.conversationUseCase.CheckShouldGenerateTitleWithProvider(provider, ctx, conversationID)
		if err != nil {
			log.Printf("Failed to check if should generate title for conversation %s: %v", conversationID, err)
			return nil
		}
		
		log.Printf("Should generate title for conversation %s: %v", conversationID, shouldGenerate)
		shouldGenerateTitle = shouldGenerate
		if shouldGenerate {
			log.Printf("First exchange complete for conversation %s, will generate title", conversationID)
			// Get the conversation messages - since we just added the assistant message, 
			// the conversation now has 2 messages total (user + assistant)
			allMessages, err := provider.Message().GetByConversationID(ctx, conversationID, 10, 0)
			if err != nil {
				log.Printf("Failed to get messages for title generation: %v", err)
				return nil
			}
			log.Printf("Retrieved %d messages for title generation", len(allMessages))
			// Find the user message (should be the first one chronologically)
			for _, msg := range allMessages {
				if msg.Role == entities.MessageRoleUser {
					userMessage = msg.Content
					log.Printf("Found user message for title generation: %s", userMessage[:min(50, len(userMessage))])
					break
				}
			}
			if userMessage == "" {
				log.Printf("No user message found for title generation")
			}
		}
		
		return nil
	})

	if err != nil {
		log.Printf("Failed to save assistant's response for conversation %s: %v", conversationID, err)
	} else if shouldGenerateTitle && userMessage != "" {
		log.Printf("Triggering title generation for conversation %s", conversationID)
		// Trigger title generation asynchronously
		uc.conversationUseCase.GenerateConversationTitle(ctx, conversationID, userMessage, responseContent.String())
	} else {
		log.Printf("Not triggering title generation - shouldGenerateTitle: %v, userMessage empty: %v", shouldGenerateTitle, userMessage == "")
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

// Note: Conversation CRUD operations have been moved to ConversationUseCase.
// This ChatUseCase now focuses on messaging and streaming functionality. 
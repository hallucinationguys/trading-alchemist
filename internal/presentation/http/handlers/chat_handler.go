package handlers

import (
	"trading-alchemist/internal/application/dto"
	"trading-alchemist/internal/application/usecases"
	"trading-alchemist/internal/presentation/responses"
	"trading-alchemist/pkg/utils"

	"bufio"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ChatHandler handles chat-related requests.
type ChatHandler struct {
	chatUseCase         *usecases.ChatUseCase
	conversationUseCase *usecases.ConversationUseCase
}

// NewChatHandler creates a new ChatHandler.
func NewChatHandler(chatUseCase *usecases.ChatUseCase, conversationUseCase *usecases.ConversationUseCase) *ChatHandler {
	return &ChatHandler{
		chatUseCase:         chatUseCase,
		conversationUseCase: conversationUseCase,
	}
}

// CreateConversation creates a new chat conversation.
// @Summary Create a new conversation
// @Description Creates a new chat session for the authenticated user.
// @Tags Chat
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateConversationRequest true "Conversation creation request"
// @Success 201 {object} responses.SuccessResponse{data=dto.ConversationDetailResponse} "Conversation created successfully"
// @Failure 400 {object} responses.ErrorResponse "Invalid request body or validation error"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /conversations [post]
func (h *ChatHandler) CreateConversation(c *fiber.Ctx) error {
	var req dto.CreateConversationRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	// Extract user from context
	userClaims, ok := c.Locals("user").(*utils.Claims)
	if !ok || userClaims == nil {
		return responses.SendError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
	}

	userID, err := uuid.Parse(userClaims.Subject)
	if err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid user ID format in token")
	}
	req.UserID = userID

	// Call conversation use case
	conversation, err := h.conversationUseCase.CreateConversation(c.Context(), &req)
	if err != nil {
		return responses.HandleError(c, err)
	}

	return responses.SendCreated(c, conversation, "Conversation created successfully")
}

// GetConversations retrieves a paginated list of conversations for the current user.
// @Summary Get user's conversations
// @Description Retrieves a paginated list of active conversations for the authenticated user.
// @Tags Chat
// @Accept json
// @Produce json
// @Security Bearer
// @Param limit query int false "Number of conversations to return" default(20)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} responses.SuccessResponse{data=[]dto.ConversationSummaryResponse} "Conversations retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /conversations [get]
func (h *ChatHandler) GetConversations(c *fiber.Ctx) error {
	// Extract user from context
	userClaims, ok := c.Locals("user").(*utils.Claims)
	if !ok || userClaims == nil {
		return responses.SendError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
	}

	userID, err := uuid.Parse(userClaims.Subject)
	if err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid user ID format in token")
	}

	// Get pagination parameters
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	// Call conversation use case
	conversations, err := h.conversationUseCase.GetUserConversations(c.Context(), userID, limit, offset)
	if err != nil {
		return responses.HandleError(c, err)
	}

	return responses.SendSuccess(c, conversations)
}

// GetConversation retrieves the details of a single conversation.
// @Summary Get conversation details
// @Description Retrieves the full details of a single conversation, including its messages, for the authenticated user.
// @Tags Chat
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Conversation ID"
// @Success 200 {object} responses.SuccessResponse{data=dto.ConversationDetailResponse} "Conversation details retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse "Invalid conversation ID format"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Forbidden - User does not own this conversation"
// @Failure 404 {object} responses.ErrorResponse "Conversation not found"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /conversations/{id} [get]
func (h *ChatHandler) GetConversation(c *fiber.Ctx) error {
	// Extract user from context
	userClaims, ok := c.Locals("user").(*utils.Claims)
	if !ok || userClaims == nil {
		return responses.SendError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
	}

	userID, err := uuid.Parse(userClaims.Subject)
	if err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid user ID format in token")
	}

	// Get conversation ID from URL
	conversationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid conversation ID format")
	}

	// Call conversation use case
	conversation, err := h.conversationUseCase.GetConversationDetails(c.Context(), conversationID, userID)
	if err != nil {
		return responses.HandleError(c, err)
	}

	return responses.SendSuccess(c, conversation)
}

// PostMessage adds a new message to a conversation and streams the response back.
// @Summary Post a message and get a streaming response
// @Description Sends a message to a conversation and streams the LLM's response back using Server-Sent Events (SSE).
// @Tags Chat
// @Accept json
// @Produce plain
// @Security Bearer
// @Param id path string true "Conversation ID"
// @Param request body dto.PostMessageRequest true "Message content"
// @Success 200 {string} string "text/event-stream response"
// @Failure 400 {object} responses.ErrorResponse "Invalid request body or ID format"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Forbidden - User does not own this conversation"
// @Failure 404 {object} responses.ErrorResponse "Conversation not found"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /conversations/{id}/messages [post]
func (h *ChatHandler) PostMessage(c *fiber.Ctx) error {
	var req dto.PostMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	userClaims, ok := c.Locals("user").(*utils.Claims)
	if !ok || userClaims == nil {
		return responses.SendError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
	}

	userID, err := uuid.Parse(userClaims.Subject)
	if err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid user ID format in token")
	}

	conversationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid conversation ID format")
	}

	// Call use case to get the stream
	eventChannel, err := h.chatUseCase.PostMessage(c.Context(), conversationID, userID, &req)
	if err != nil {
		return responses.HandleError(c, err)
	}

	// Set headers for SSE
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		for event := range eventChannel {
			if event.Error != nil {
				log.Printf("SSE stream error for conversation %s: %v", conversationID, event.Error)
				// Optionally, you could send an error event to the client
				// For now, we just close the connection by returning.
				return
			}

			// Marshal the event to JSON
			jsonEvent, err := json.Marshal(event)
			if err != nil {
				log.Printf("Error marshaling stream event: %v", err)
				continue // Skip this event
			}

			// Write the event in SSE format
			if _, err := fmt.Fprintf(w, "data: %s\n\n", jsonEvent); err != nil {
				log.Printf("Error writing to SSE stream: %v", err)
				return // Stop streaming if we can't write
			}

			// Flush the writer to send the event immediately
			if err := w.Flush(); err != nil {
				log.Printf("Error flushing SSE stream: %v", err)
				return // Stop streaming if we can't flush
			}
		}
	})

	return nil
}

// GetAvailableTools retrieves a list of available tools.
// @Summary Get available tools
// @Description Retrieves a list of all active tools that can be used by the LLM. Can be filtered by provider.
// @Tags Chat
// @Accept json
// @Produce json
// @Security Bearer
// @Param provider_id query string false "Filter tools by a specific provider ID"
// @Success 200 {object} responses.SuccessResponse{data=[]dto.ToolResponse} "Tools retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /tools [get]
func (h *ChatHandler) GetAvailableTools(c *fiber.Ctx) error {
	// We could optionally filter by provider ID from the query string
	var providerID *uuid.UUID
	if providerIDStr := c.Query("provider_id"); providerIDStr != "" {
		if id, err := uuid.Parse(providerIDStr); err == nil {
			providerID = &id
		}
	}

	tools, err := h.chatUseCase.GetAvailableTools(c.Context(), providerID)
	if err != nil {
		return responses.HandleError(c, err)
	}

	return responses.SendSuccess(c, tools)
}

// UpdateConversationTitle updates the title of a conversation.
// @Summary Update conversation title
// @Description Updates the title of a conversation for the authenticated user.
// @Tags Chat
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Conversation ID"
// @Param request body dto.UpdateConversationTitleRequest true "Title update request"
// @Success 200 {object} responses.SuccessResponse "Conversation title updated successfully"
// @Failure 400 {object} responses.ErrorResponse "Invalid request body or ID format"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Forbidden - User does not own this conversation"
// @Failure 404 {object} responses.ErrorResponse "Conversation not found"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /conversations/{id}/title [put]
func (h *ChatHandler) UpdateConversationTitle(c *fiber.Ctx) error {
	var req dto.UpdateConversationTitleRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	// Extract user from context
	userClaims, ok := c.Locals("user").(*utils.Claims)
	if !ok || userClaims == nil {
		return responses.SendError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
	}

	userID, err := uuid.Parse(userClaims.Subject)
	if err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid user ID format in token")
	}

	// Get conversation ID from URL
	conversationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid conversation ID format")
	}

	// Call conversation use case
	err = h.conversationUseCase.UpdateConversationTitle(c.Context(), conversationID, userID, &req)
	if err != nil {
		return responses.HandleError(c, err)
	}

	return responses.SendSuccess(c, nil)
}

// ArchiveConversation archives (soft deletes) a conversation.
// @Summary Archive conversation
// @Description Archives a conversation for the authenticated user.
// @Tags Chat
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Conversation ID"
// @Success 200 {object} responses.SuccessResponse "Conversation archived successfully"
// @Failure 400 {object} responses.ErrorResponse "Invalid conversation ID format"
// @Failure 401 {object} responses.ErrorResponse "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse "Forbidden - User does not own this conversation"
// @Failure 404 {object} responses.ErrorResponse "Conversation not found"
// @Failure 500 {object} responses.ErrorResponse "Internal server error"
// @Router /conversations/{id} [delete]
func (h *ChatHandler) ArchiveConversation(c *fiber.Ctx) error {
	// Extract user from context
	userClaims, ok := c.Locals("user").(*utils.Claims)
	if !ok || userClaims == nil {
		return responses.SendError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
	}

	userID, err := uuid.Parse(userClaims.Subject)
	if err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid user ID format in token")
	}

	// Get conversation ID from URL
	conversationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return responses.SendError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid conversation ID format")
	}

	// Call conversation use case
	err = h.conversationUseCase.ArchiveConversation(c.Context(), conversationID, userID)
	if err != nil {
		return responses.HandleError(c, err)
	}

	return responses.SendSuccess(c, nil)
}

// TODO: Implement handler methods:
// - PostMessage(c *fiber.Ctx) error 
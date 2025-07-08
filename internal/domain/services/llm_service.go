package services

import (
	"context"
	"trading-alchemist/internal/domain/chat"
)

// ChatStreamEvent represents a single event in a chat completion stream.
type ChatStreamEvent struct {
	ContentDelta string      `json:"content_delta"`
	ToolCall     interface{} `json:"tool_call,omitempty"` // Placeholder for tool call data
	IsLast       bool        `json:"is_last"`
	Error        error       `json:"error,omitempty"`
}

// LLMService defines the interface for interacting with a Large Language Model.
type LLMService interface {
	// StreamChatCompletion sends a chat request and streams the response.
	StreamChatCompletion(
		ctx context.Context,
		provider *chat.Provider,
		model *chat.Model,
		messages []*chat.Message,
		apiKey string,
		apiBaseOverride string,
	) (<-chan ChatStreamEvent, error)
} 
package chat

import (
	"context"

	"github.com/google/uuid"
)

type ToolRepository interface {
	GetAvailableTools(ctx context.Context, providerID *uuid.UUID) ([]*Tool, error)
	LogToolUsage(ctx context.Context, messageTool *MessageTool) error
} 
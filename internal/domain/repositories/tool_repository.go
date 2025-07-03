package repositories

import (
	"context"
	"trading-alchemist/internal/domain/entities"

	"github.com/google/uuid"
)

type ToolRepository interface {
	GetAvailableTools(ctx context.Context, providerID *uuid.UUID) ([]*entities.Tool, error)
	LogToolUsage(ctx context.Context, messageTool *entities.MessageTool) error
} 
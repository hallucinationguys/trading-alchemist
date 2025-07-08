package chat

import (
	"time"
	"trading-alchemist/internal/domain/shared"

	"github.com/google/uuid"
)

// MessageTool represents tool usage in messages
type MessageTool struct {
	ID         uuid.UUID `json:"id" db:"id"`
	MessageID  uuid.UUID `json:"message_id" db:"message_id"`
	ToolID     uuid.UUID `json:"tool_id" db:"tool_id"`
	Input      shared.JSONB     `json:"input" db:"input"`
	Output     shared.JSONB     `json:"output" db:"output"`
	ExecutedAt time.Time `json:"executed_at" db:"executed_at"`
	Duration   int64     `json:"duration" db:"duration"` // Execution time in ms
	Success    bool      `json:"success" db:"success"`
	Error      *string   `json:"error" db:"error"`
} 
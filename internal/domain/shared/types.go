package shared

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONB represents a JSONB database type.
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface.
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface.
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleSystem    MessageRole = "system"
	MessageRoleTool      MessageRole = "tool"
)

type ArtifactType string

const (
	ArtifactTypeCode     ArtifactType = "code"
	ArtifactTypeDocument ArtifactType = "document"
	ArtifactTypeChart    ArtifactType = "chart"
	ArtifactTypeImage    ArtifactType = "image"
	ArtifactTypeHTML     ArtifactType = "html"
	ArtifactTypeSVG      ArtifactType = "svg"
) 
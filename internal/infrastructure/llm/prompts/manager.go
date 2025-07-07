package prompts

import (
	"bytes"
	"fmt"
	"text/template"
)

// Prompt represents a structured prompt with metadata
type Prompt struct {
	Name        string
	Description string
	Version     string
	SystemPrompt string
	UserTemplate string
}

// PromptManager handles loading and processing prompts
type PromptManager struct {
	prompts map[string]*Prompt
}

// NewPromptManager creates a new prompt manager instance
func NewPromptManager() *PromptManager {
	pm := &PromptManager{
		prompts: make(map[string]*Prompt),
	}
	
	// Register all available prompts
	pm.registerPrompts()
	
	return pm
}

// registerPrompts registers all available prompts in the system
func (pm *PromptManager) registerPrompts() {
	// Register title generation prompt
	pm.prompts["title_generation"] = GetTitleGenerationPrompt()
	
	// Future prompts can be registered here:
	// pm.prompts["code_analysis"] = GetCodeAnalysisPrompt()
	// pm.prompts["creative_writing"] = GetCreativeWritingPrompt()
}

// GetPrompt retrieves a prompt by name
func (pm *PromptManager) GetPrompt(name string) (*Prompt, error) {
	prompt, exists := pm.prompts[name]
	if !exists {
		return nil, fmt.Errorf("prompt not found: %s", name)
	}
	return prompt, nil
}

// RenderUserPrompt renders the user prompt template with the provided data
func (pm *PromptManager) RenderUserPrompt(promptName string, data interface{}) (string, error) {
	prompt, err := pm.GetPrompt(promptName)
	if err != nil {
		return "", err
	}
	
	tmpl, err := template.New(promptName).Parse(prompt.UserTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template for %s: %w", promptName, err)
	}
	
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template for %s: %w", promptName, err)
	}
	
	return buf.String(), nil
}

// GetSystemPrompt returns the system prompt for a given prompt name
func (pm *PromptManager) GetSystemPrompt(promptName string) (string, error) {
	prompt, err := pm.GetPrompt(promptName)
	if err != nil {
		return "", err
	}
	return prompt.SystemPrompt, nil
}

// ListPrompts returns a list of all available prompt names
func (pm *PromptManager) ListPrompts() []string {
	names := make([]string, 0, len(pm.prompts))
	for name := range pm.prompts {
		names = append(names, name)
	}
	return names
} 
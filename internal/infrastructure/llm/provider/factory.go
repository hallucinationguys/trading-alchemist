package provider

import (
	"fmt"
)

// NewClientForProvider creates a provider-specific client with the given API key.
func NewClientForProvider(providerName, apiKey, apiBaseOverride string) (ProviderClient, error) {
	switch providerName {
	case "openai":
		return NewOpenAIClient(apiKey, apiBaseOverride)
	// case "google":
	// 	return NewGoogleClient(apiKey, apiBaseOverride)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerName)
	}
} 
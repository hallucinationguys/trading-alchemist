package email

import (
	"fmt"

	"trading-alchemist/internal/config"
	"trading-alchemist/internal/domain/services"
)

// NewEmailService creates an email service based on configuration
func NewEmailService(cfg *config.Config) (services.EmailService, error) {
	// Only use Resend provider
	if cfg.Email.ResendAPIKey != "" {
		return NewResendProvider(cfg), nil
	}

	return nil, fmt.Errorf("no email provider configured: please set RESEND_API_KEY")
} 
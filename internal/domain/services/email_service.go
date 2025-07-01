package services

import (
	"context"

	"trading-alchemist/internal/domain/entities"
)

type EmailService interface {
	// SendMagicLinkEmail sends a magic link email to the user
	SendMagicLinkEmail(ctx context.Context, user *entities.User, magicLink *entities.MagicLink) error
	
	// SendWelcomeEmail sends a welcome email to new users
	SendWelcomeEmail(ctx context.Context, user *entities.User) error
	
	// SendEmailVerificationEmail sends an email verification email
	SendEmailVerificationEmail(ctx context.Context, user *entities.User, magicLink *entities.MagicLink) error
} 
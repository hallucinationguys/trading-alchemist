package usecases

import (
	"context"
	"fmt"
	"net"
	"time"

	"trading-alchemist/internal/application/dto"
	"trading-alchemist/internal/config"
	"trading-alchemist/internal/domain/entities"
	"trading-alchemist/internal/domain/repositories"
	"trading-alchemist/internal/domain/services"
	"trading-alchemist/pkg/errors"
	"trading-alchemist/pkg/utils"

	"github.com/google/uuid"
)

type AuthUseCase struct {
	userRepo      repositories.UserRepository
	magicLinkRepo repositories.MagicLinkRepository
	emailService  services.EmailService
	tokenService  services.TokenService
	config        *config.Config
}

// NewAuthUseCase creates a new authentication use case
func NewAuthUseCase(
	userRepo repositories.UserRepository,
	magicLinkRepo repositories.MagicLinkRepository,
	emailService services.EmailService,
	tokenService services.TokenService,
	config *config.Config,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:      userRepo,
		magicLinkRepo: magicLinkRepo,
		emailService:  emailService,
		tokenService:  tokenService,
		config:        config,
	}
}

// SendMagicLink sends a magic link to the user's email
func (uc *AuthUseCase) SendMagicLink(ctx context.Context, req *dto.SendMagicLinkRequest) (*dto.SendMagicLinkResponse, error) {
	// Validate email
	if !utils.IsValidEmail(req.Email) {
		return nil, errors.NewAppError(errors.CodeValidation, "Invalid email address", errors.ErrInvalidEmail)
	}

	// Normalize email
	email := utils.NormalizeEmail(req.Email)

	// Set default purpose if not provided
	purpose := entities.MagicLinkPurposeLogin
	if req.Purpose != "" {
		purpose = entities.MagicLinkPurpose(req.Purpose)
	}

	// Get or create user
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Create new user if not found
		if err == errors.ErrUserNotFound {
			newUser := &entities.User{
				ID:            uuid.New(),
				Email:         email,
				EmailVerified: false,
				IsActive:      true,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			user, err = uc.userRepo.Create(ctx, newUser)
			if err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}

			// Send welcome email for new users
			go func() {
				if err := uc.emailService.SendWelcomeEmail(context.Background(), user); err != nil {
					// Log error but don't fail the main flow
					fmt.Printf("Failed to send welcome email: %v\n", err)
				}
			}()
		} else {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
	}

	// Invalidate existing magic links for this purpose
	if err := uc.magicLinkRepo.InvalidateUserLinks(ctx, user.ID, purpose); err != nil {
		return nil, fmt.Errorf("failed to invalidate existing links: %w", err)
	}

	// Generate new magic link
	token, err := uc.tokenService.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	tokenHash := uc.tokenService.HashToken(token)
	
	// Parse magic link TTL
	magicLinkTTL, err := time.ParseDuration(uc.config.App.MagicLinkTTL)
	if err != nil {
		return nil, fmt.Errorf("invalid MAGIC_LINK_TTL configuration: %w", err)
	}
	expiresAt := time.Now().Add(magicLinkTTL)

	magicLink := &entities.MagicLink{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		Purpose:   purpose,
		CreatedAt: time.Now(),
	}

	// Set IP address and user agent if provided
	if req.IPAddress != "" {
		if ip := net.ParseIP(req.IPAddress); ip != nil {
			magicLink.IPAddress = &req.IPAddress
		}
	}
	if req.UserAgent != "" {
		magicLink.UserAgent = &req.UserAgent
	}

	// Save magic link
	createdLink, err := uc.magicLinkRepo.Create(ctx, magicLink)
	if err != nil {
		return nil, fmt.Errorf("failed to create magic link: %w", err)
	}

	// Send email
	var emailErr error
	switch purpose {
	case entities.MagicLinkPurposeEmailVerification:
		emailErr = uc.emailService.SendEmailVerificationEmail(ctx, user, createdLink)
	default:
		emailErr = uc.emailService.SendMagicLinkEmail(ctx, user, createdLink)
	}

	if emailErr != nil {
		return nil, fmt.Errorf("failed to send email: %w", emailErr)
	}

	return &dto.SendMagicLinkResponse{
		Message: "Magic link sent successfully",
		Sent:    true,
	}, nil
}

// VerifyMagicLink verifies a magic link and returns authentication tokens
func (uc *AuthUseCase) VerifyMagicLink(ctx context.Context, req *dto.VerifyMagicLinkRequest) (*dto.VerifyMagicLinkResponse, error) {
	// Get magic link and user
	magicLink, user, err := uc.magicLinkRepo.GetByToken(ctx, req.Token)
	if err != nil {
		if err == errors.ErrMagicLinkNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "Invalid or expired magic link", err)
		}
		return nil, fmt.Errorf("failed to get magic link: %w", err)
	}

	// Verify token hash
	if !uc.tokenService.VerifyToken(req.Token, magicLink.TokenHash) {
		return nil, errors.NewAppError(errors.CodeUnauthorized, "Invalid magic link", errors.ErrInvalidToken)
	}

	// Check if link is valid
	if !magicLink.IsValid() {
		if magicLink.IsExpired() {
			return nil, errors.NewAppError(errors.CodeUnauthorized, "Magic link has expired", errors.ErrMagicLinkExpired)
		}
		if magicLink.IsUsed() {
			return nil, errors.NewAppError(errors.CodeUnauthorized, "Magic link has already been used", errors.ErrMagicLinkAlreadyUsed)
		}
	}

	// Mark magic link as used
	if _, err := uc.magicLinkRepo.MarkAsUsed(ctx, magicLink.ID); err != nil {
		return nil, fmt.Errorf("failed to mark magic link as used: %w", err)
	}

	// Handle email verification
	if magicLink.Purpose == entities.MagicLinkPurposeEmailVerification && !user.EmailVerified {
		if _, err := uc.userRepo.VerifyEmail(ctx, user.ID); err != nil {
			return nil, fmt.Errorf("failed to verify user email: %w", err)
		}
		user.EmailVerified = true
	}

	// Parse JWT TTL
	jwtTTL, err := time.ParseDuration(uc.config.JWT.TTL)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_TTL configuration: %w", err)
	}

	// Generate JWT token
	token, err := uc.tokenService.GenerateJWT(user.ID, jwtTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	return &dto.VerifyMagicLinkResponse{
		User:        dto.ToUserResponse(user),
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(jwtTTL.Seconds()),
	}, nil
}

// ValidateToken validates a JWT token and returns the user
func (uc *AuthUseCase) ValidateToken(ctx context.Context, token string) (*entities.User, error) {
	userID, err := uc.tokenService.ValidateJWT(token)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeUnauthorized, "Invalid token", err)
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == errors.ErrUserNotFound {
			return nil, errors.NewAppError(errors.CodeUnauthorized, "User not found", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.IsAccountActive() {
		return nil, errors.NewAppError(errors.CodeForbidden, "Account is inactive", errors.ErrForbidden)
	}

	return user, nil
} 
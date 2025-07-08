package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"trading-alchemist/internal/config"
	"trading-alchemist/internal/domain/auth"
	"trading-alchemist/internal/domain/services"
	"trading-alchemist/internal/infrastructure/database"
	"trading-alchemist/pkg/errors"
	"trading-alchemist/pkg/utils"

	"github.com/google/uuid"
)

type AuthUseCase struct {
	emailService  services.EmailService
	config        *config.Config
	dbService     *database.Service
}

// NewAuthUseCase creates a new authentication use case
func NewAuthUseCase(
	emailService services.EmailService,
	config *config.Config,
	dbService *database.Service,
) *AuthUseCase {
	return &AuthUseCase{
		emailService:  emailService,
		config:        config,
		dbService:     dbService,
	}
}

// SendMagicLink sends a magic link to the user's email
func (uc *AuthUseCase) SendMagicLink(ctx context.Context, req *SendMagicLinkRequest) (*SendMagicLinkResponse, error) {
	if !utils.IsValidEmail(req.Email) {
		return nil, errors.NewAppError(errors.CodeValidation, "Invalid email address", errors.ErrInvalidEmail)
	}

	email := utils.NormalizeEmail(req.Email)
	purpose := auth.MagicLinkPurposeLogin
	if req.Purpose != "" {
		purpose = auth.MagicLinkPurpose(req.Purpose)
	}

	var user *auth.User
	var createdLink *auth.MagicLink

	err := uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		var err error
		userRepo := provider.User()
		magicLinkRepo := provider.MagicLink()

		user, err = userRepo.GetByEmail(ctx, email)
		if err != nil {
			if err == errors.ErrUserNotFound {
				newUser := &auth.User{
					ID:            uuid.New(),
					Email:         email,
					EmailVerified: false,
					IsActive:      true,
				}
				user, err = userRepo.Create(ctx, newUser)
				if err != nil {
					return fmt.Errorf("failed to create user: %w", err)
				}
				go func(userToSend *auth.User) {
					if err := uc.emailService.SendWelcomeEmail(context.Background(), userToSend); err != nil {
						log.Printf("Failed to send welcome email: %v\n", err)
					}
				}(user)
			} else {
				return fmt.Errorf("failed to get user: %w", err)
			}
		}

		if err := magicLinkRepo.InvalidateUserLinks(ctx, user.ID, purpose); err != nil {
			return fmt.Errorf("failed to invalidate existing links: %w", err)
		}

		token, err := utils.GenerateSecureToken(32)
		if err != nil {
			return fmt.Errorf("failed to generate token: %w", err)
		}

		magicLinkTTL, err := time.ParseDuration(uc.config.App.MagicLinkTTL)
		if err != nil {
			return fmt.Errorf("invalid MAGIC_LINK_TTL configuration: %w", err)
		}

		magicLink := &auth.MagicLink{
			ID:        uuid.New(),
			UserID:    user.ID,
			Token:     token,
			TokenHash: utils.HashToken(token),
			ExpiresAt: time.Now().Add(magicLinkTTL),
			Purpose:   purpose,
		}
		if req.IPAddress != "" {
			magicLink.IPAddress = &req.IPAddress
		}
		if req.UserAgent != "" {
			magicLink.UserAgent = &req.UserAgent
		}

		createdLink, err = magicLinkRepo.Create(ctx, magicLink)
		if err != nil {
			return fmt.Errorf("failed to create magic link: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err // The error from ExecuteInTx is already descriptive
	}

	// Send email asynchronously after transaction is committed
	go func(userToSend *auth.User, linkToSend *auth.MagicLink) {
		var emailErr error
		switch purpose {
		case auth.MagicLinkPurposeEmailVerification:
			emailErr = uc.emailService.SendEmailVerificationEmail(context.Background(), userToSend, linkToSend)
		default:
			emailErr = uc.emailService.SendMagicLinkEmail(context.Background(), userToSend, linkToSend)
		}
		if emailErr != nil {
			log.Printf("Failed to send magic link email asynchronously: %v\n", emailErr)
		}
	}(user, createdLink)

	return &SendMagicLinkResponse{
		Message: "Magic link sent successfully",
		Sent:    true,
	}, nil
}

// VerifyMagicLink verifies a magic link and returns authentication tokens
func (uc *AuthUseCase) VerifyMagicLink(ctx context.Context, req *VerifyMagicLinkRequest) (*VerifyMagicLinkResponse, error) {
	var magicLink *auth.MagicLink
	var user *auth.User
	
	err := uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		var err error
		magicLink, user, err = provider.MagicLink().GetByToken(ctx, req.Token)
		return err
	})
	
	if err != nil {
		if err == errors.ErrMagicLinkNotFound {
			return nil, errors.NewAppError(errors.CodeNotFound, "Invalid or expired magic link", err)
		}
		return nil, fmt.Errorf("failed to get magic link: %w", err)
	}

	if !utils.VerifyTokenHash(req.Token, magicLink.TokenHash) {
		return nil, errors.NewAppError(errors.CodeUnauthorized, "Invalid magic link", errors.ErrInvalidToken)
	}

	if !magicLink.IsValid() {
		if magicLink.IsExpired() {
			return nil, errors.NewAppError(errors.CodeUnauthorized, "Magic link has expired", errors.ErrMagicLinkExpired)
		}
		if magicLink.IsUsed() {
			return nil, errors.NewAppError(errors.CodeUnauthorized, "Magic link has already been used", errors.ErrMagicLinkAlreadyUsed)
		}
	}

	err = uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		if _, err := provider.MagicLink().MarkAsUsed(ctx, magicLink.ID); err != nil {
			return fmt.Errorf("failed to mark magic link as used: %w", err)
		}
		if magicLink.Purpose == auth.MagicLinkPurposeEmailVerification && !user.EmailVerified {
			if _, err := provider.User().VerifyEmail(ctx, user.ID); err != nil {
				return fmt.Errorf("failed to verify user email: %w", err)
			}
			user.EmailVerified = true
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	jwtTTL, err := time.ParseDuration(uc.config.JWT.TTL)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_TTL configuration: %w", err)
	}

	token, err := utils.GenerateJWT(user, uc.config.JWT.Secret, jwtTTL, uc.config.App.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	return &VerifyMagicLinkResponse{
		User:        ToUserResponse(user),
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(jwtTTL.Seconds()),
	}, nil
}

// ValidateToken validates a JWT token and returns the user claims
func (uc *AuthUseCase) ValidateToken(ctx context.Context, token string) (*utils.Claims, error) {
	// Validate JWT using utility function
	claims, err := utils.ValidateJWT(token, uc.config.JWT.Secret)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeUnauthorized, "Invalid token", err)
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeUnauthorized, "Invalid user ID in token", err)
	}

	var user *auth.User
	err = uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		var err error
		user, err = provider.User().GetByID(ctx, userID)
		if err != nil {
			if err == errors.ErrUserNotFound {
				return errors.NewAppError(errors.CodeUnauthorized, "User not found", err)
			}
			return fmt.Errorf("failed to get user: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if !user.IsAccountActive() {
		return nil, errors.NewAppError(errors.CodeForbidden, "Account is inactive", errors.ErrForbidden)
	}

	// The email in claims might be stale, so we prefer the one from the database.
	claims.Email = user.Email
	return claims, nil
} 
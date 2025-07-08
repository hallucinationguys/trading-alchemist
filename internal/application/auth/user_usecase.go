package auth

import (
	"context"
	"fmt"
	"trading-alchemist/internal/domain/auth"
	"trading-alchemist/internal/infrastructure/database"

	"github.com/google/uuid"
)

type UserUseCase struct {
	dbService *database.Service
}

func NewUserUseCase(dbService *database.Service) *UserUseCase {
	return &UserUseCase{
		dbService: dbService,
	}
}

func (uc *UserUseCase) GetUserProfile(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	var user *auth.User
	err := uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		var err error
		user, err = provider.User().GetByID(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user profile: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	userResponse := ToUserResponse(user)
	return &userResponse, nil
}

// UpdateUserProfile updates the user profile information
func (uc *UserUseCase) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *UpdateUserRequest) (*UserResponse, error) {
	var updatedUser *auth.User
	
	err := uc.dbService.ExecuteInTx(ctx, func(provider database.RepositoryProvider) error {
		// First, get the current user to ensure they exist
		user, err := provider.User().GetByID(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user for update: %w", err)
		}

		// Update the fields that were provided in the request
		if req.FirstName != nil {
			user.FirstName = req.FirstName
		}
		if req.LastName != nil {
			user.LastName = req.LastName
		}
		if req.AvatarURL != nil {
			user.AvatarURL = req.AvatarURL
		}

		// Update the user in the database
		updatedUser, err = provider.User().Update(ctx, user)
		if err != nil {
			return fmt.Errorf("failed to update user profile: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	userResponse := ToUserResponse(updatedUser)
	return &userResponse, nil
} 
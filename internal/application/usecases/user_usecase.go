package usecases

import (
	"context"
	"fmt"
	"trading-alchemist/internal/application/dto"
	"trading-alchemist/internal/domain/repositories"

	"github.com/google/uuid"
)

type UserUseCase struct {
	userRepo repositories.UserRepository
}

func NewUserUseCase(userRepo repositories.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UserUseCase) GetUserProfile(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	userResponse := dto.ToUserResponse(user)
	return &userResponse, nil
}

// UpdateUserProfile updates the user profile information
func (uc *UserUseCase) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// First, get the current user to ensure they exist
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user for update: %w", err)
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
	updatedUser, err := uc.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	userResponse := dto.ToUserResponse(updatedUser)
	return &userResponse, nil
} 
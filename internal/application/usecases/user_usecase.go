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
package services

import (
	"time"

	"github.com/google/uuid"
)

type TokenService interface {
	// GenerateToken generates a secure random token
	GenerateToken() (string, error)
	
	// HashToken hashes a token for secure storage
	HashToken(token string) string
	
	// VerifyToken verifies if a plain token matches its hash
	VerifyToken(plainToken, hashedToken string) bool
	
	// GenerateJWT generates a JWT token for a user
	GenerateJWT(userID uuid.UUID, expiresIn time.Duration) (string, error)
	
	// ValidateJWT validates a JWT token and returns the user ID
	ValidateJWT(token string) (uuid.UUID, error)
} 
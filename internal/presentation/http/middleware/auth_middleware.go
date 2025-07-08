package middleware

import (
	"strings"
	"trading-alchemist/internal/application/auth"
	"trading-alchemist/internal/presentation/responses"

	"github.com/gofiber/fiber/v2"
)

func NewAuthMiddleware(authUseCase *auth.AuthUseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return responses.SendError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return responses.SendError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization header format")
		}

		token := parts[1]
		claims, err := authUseCase.ValidateToken(c.Context(), token)
		if err != nil {
			return responses.HandleError(c, err)
		}

		c.Locals("user", claims)
		return c.Next()
	}
} 
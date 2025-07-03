package routes

import (
	"github.com/gofiber/fiber/v2"

	"trading-alchemist/internal/application/usecases"
	"trading-alchemist/internal/config"
	"trading-alchemist/internal/presentation/http/handlers"
	"trading-alchemist/internal/presentation/http/middleware"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, cfg *config.Config, authUseCase *usecases.AuthUseCase, userUseCase *usecases.UserUseCase, chatUseCase *usecases.ChatUseCase, providerUseCase *usecases.UserProviderSettingUseCase, modelAvailabilityUseCase *usecases.ModelAvailabilityUseCase) {
	// Create handlers
	authHandler := handlers.NewAuthHandler(authUseCase)
	userHandler := handlers.NewUserHandler(userUseCase, authUseCase)
	chatHandler := handlers.NewChatHandler(chatUseCase)
	providerHandler := handlers.NewProviderHandler(providerUseCase, modelAvailabilityUseCase)

	// Create auth middleware
	authMiddleware := middleware.NewAuthMiddleware(authUseCase)

	// Health check endpoints (outside API versioning for monitoring)
	app.Get("/health", handlers.CheckHealth)

	// API Documentation endpoints
	setupDocumentationRoutes(app)

	// API versioning - v1
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Setup routes for each handler
	setupHealthRoutes(v1)
	setupV1AuthRoutes(v1, authHandler)
	setupV1UserRoutes(v1, userHandler, authMiddleware)
	setupV1ChatRoutes(v1, chatHandler, authMiddleware)
	setupV1ProviderRoutes(v1, providerHandler, authMiddleware)
}

// setupDocumentationRoutes sets up Swagger documentation routes
func setupDocumentationRoutes(app *fiber.App) {
	// Swagger JSON endpoint
	app.Get("/swagger.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})
	
	// Swagger UI redirect
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("https://petstore.swagger.io/?url=" + c.BaseURL() + "/swagger.json")
	})

	// API information endpoint
	app.Get("/api", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":        "Trading Alchemist API",
			"version":     "1.0.0",
			"description": "RESTful API for Trading Alchemist platform",
			"versions": fiber.Map{
				"v1": fiber.Map{
					"status":      "active",
					"description": "Current stable version",
				},
			},
			"documentation": fiber.Map{
				"swagger_ui": "/docs",
				"openapi":    "/swagger.json",
			},
			"health_checks": []string{
				"/health",
			},
		})
	})
}

// setupHealthRoutes configures health check routes
func setupHealthRoutes(v1 fiber.Router) {
	health := v1.Group("/health")
	health.Get("", handlers.CheckHealth)
}

// setupV1AuthRoutes configures v1 authentication routes
func setupV1AuthRoutes(v1 fiber.Router, authHandler *handlers.AuthHandler) {
	auth := v1.Group("/auth")

	// Authentication endpoints
	auth.Post("/magic-link", authHandler.SendMagicLink)      // POST /api/v1/auth/magic-link

	// NOTE: This endpoint is POST to allow the frontend to securely send
	// the token in the request body after extracting it from the URL.
	auth.Post("/verify", authHandler.VerifyMagicLink)        // POST /api/v1/auth/verify

}

// setupV1UserRoutes configures v1 user routes
func setupV1UserRoutes(v1 fiber.Router, userHandler *handlers.UserHandler, authMiddleware fiber.Handler) {
	users := v1.Group("/users")

	// Protected user routes
	users.Use(authMiddleware)
	users.Get("/profile", userHandler.GetProfile)
	users.Put("/profile", userHandler.UpdateProfile)
}

// setupV1ChatRoutes configures v1 chat routes
func setupV1ChatRoutes(v1 fiber.Router, chatHandler *handlers.ChatHandler, authMiddleware fiber.Handler) {
	conversations := v1.Group("/conversations")
	conversations.Use(authMiddleware)

	conversations.Get("/", chatHandler.GetConversations)
	conversations.Post("/", chatHandler.CreateConversation)
	conversations.Get("/:id", chatHandler.GetConversation)
	conversations.Post("/:id/messages", chatHandler.PostMessage)

	// Tool routes
	tools := v1.Group("/tools")
	tools.Use(authMiddleware)
	tools.Get("/", chatHandler.GetAvailableTools)
}

// setupV1ProviderRoutes configures v1 provider routes
func setupV1ProviderRoutes(v1 fiber.Router, providerHandler *handlers.ProviderHandler, authMiddleware fiber.Handler) {
	providers := v1.Group("/providers")
	providers.Use(authMiddleware)

	providers.Get("/", providerHandler.ListProviders)
	providers.Get("/settings", providerHandler.ListUserSettings)
	providers.Post("/settings", providerHandler.UpsertUserSetting)
}


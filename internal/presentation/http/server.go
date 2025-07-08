package server

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"trading-alchemist/internal/application/auth"
	"trading-alchemist/internal/application/chat"
	"trading-alchemist/internal/config"
	"trading-alchemist/internal/domain/services"
	"trading-alchemist/internal/infrastructure/database"
	infraServices "trading-alchemist/internal/infrastructure/services"
	"trading-alchemist/internal/presentation/http/routes"
	"trading-alchemist/internal/presentation/responses"
)

// Server represents the HTTP server
type Server struct {
	app    *fiber.App
	config *config.Config
}

// NewServer creates a new HTTP server with all dependencies
func NewServer(cfg *config.Config, authUseCase *auth.AuthUseCase, dbService *database.Service, llmService services.LLMService) *Server {
	// Create Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		StrictRouting:  false,
		CaseSensitive:  false,
		ServerHeader:   "Trading Alchemist",
		AppName:        cfg.App.Name,
		ErrorHandler:   errorHandler,
	})

	// Add global middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Create use cases
	userUseCase := auth.NewUserUseCase(dbService)
	conversationUseCase := chat.NewConversationUseCase(dbService, cfg, llmService)
	chatUseCase := chat.NewChatUseCase(dbService, cfg, llmService, conversationUseCase)
	providerUseCase := chat.NewUserProviderSettingUseCase(dbService, cfg)
	
	// Create API key service and model availability use case
	// We create a temporary repository provider to access the user provider setting repository
	var modelAvailabilityUseCase *chat.ModelAvailabilityUseCase
	err := dbService.ExecuteInTx(context.Background(), func(provider database.RepositoryProvider) error {
		apiKeyService := infraServices.NewAPIKeyService(provider.UserProviderSetting())
		modelAvailabilityUseCase = chat.NewModelAvailabilityUseCase(dbService, apiKeyService)
		return nil
	})
	if err != nil {
		panic("Failed to create model availability use case: " + err.Error())
	}

	// Setup all routes with use cases
	routes.SetupRoutes(app, cfg, authUseCase, userUseCase, chatUseCase, conversationUseCase, providerUseCase, modelAvailabilityUseCase)

	return &Server{
		app:    app,
		config: cfg,
	}
}

// Start starts the HTTP server
func (s *Server) Start(address string) error {
	return s.app.Listen(address)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}

// errorHandler handles Fiber errors
func errorHandler(c *fiber.Ctx, err error) error {
	// Default 500 status code
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Send custom error page or JSON response
	return responses.SendError(c, code, "INTERNAL_ERROR", err.Error())
} 
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"trading-alchemist/internal/config"
	"trading-alchemist/internal/domain/services"
	"trading-alchemist/internal/presentation/http/handlers"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, cfg *config.Config, db *pgxpool.Pool, emailService services.EmailService) {
	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler()
	userHandler := handlers.NewUserHandler()

	// Health check endpoints (outside API versioning for monitoring)
	app.Get("/health", healthHandler.CheckHealth)

	// API Documentation endpoints
	setupDocumentationRoutes(app)

	// API versioning - v1
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Setup versioned routes
	setupV1AuthRoutes(v1, authHandler)
	setupV1UserRoutes(v1, userHandler)
	setupV1HealthRoutes(v1, healthHandler)
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

// setupV1AuthRoutes configures v1 authentication routes
func setupV1AuthRoutes(v1 fiber.Router, authHandler *handlers.AuthHandler) {
	auth := v1.Group("/auth")

	// Authentication endpoints
	auth.Post("/magic-link", authHandler.SendMagicLink)      // POST /api/v1/auth/magic-link
	auth.Post("/verify", authHandler.VerifyMagicLink)        // POST /api/v1/auth/verify
	
	// TODO: Implement additional auth endpoints in the future
	// auth.Post("/logout", authHandler.Logout)              // POST /api/v1/auth/logout
	// auth.Post("/refresh", authHandler.RefreshToken)       // POST /api/v1/auth/refresh
	// auth.Get("/me", authHandler.GetCurrentUser)           // GET /api/v1/auth/me
	// auth.Get("/sessions", authHandler.GetUserSessions)    // GET /api/v1/auth/sessions
}

// setupV1UserRoutes configures v1 user routes
func setupV1UserRoutes(v1 fiber.Router, userHandler *handlers.UserHandler) {
	users := v1.Group("/users")

	// User profile endpoints
	users.Get("/profile", userHandler.GetProfile)           // GET /api/v1/users/profile
	users.Put("/profile", userHandler.UpdateProfile)        // PUT /api/v1/users/profile
	
	// TODO: Implement additional user endpoints in the future
	// users.Patch("/profile", userHandler.PatchProfile)     // PATCH /api/v1/users/profile
	// users.Delete("/profile", userHandler.DeleteProfile)   // DELETE /api/v1/users/profile
	// users.Get("/settings", userHandler.GetSettings)       // GET /api/v1/users/settings
	// users.Put("/settings", userHandler.UpdateSettings)    // PUT /api/v1/users/settings
	// users.Get("/preferences", userHandler.GetPreferences) // GET /api/v1/users/preferences
	// users.Put("/preferences", userHandler.UpdatePreferences) // PUT /api/v1/users/preferences
}

// setupV1HealthRoutes configures v1 health check routes
func setupV1HealthRoutes(v1 fiber.Router, healthHandler *handlers.HealthHandler) {
	health := v1.Group("/health")

	// Comprehensive health checks
	health.Get("/", healthHandler.CheckHealth)               // GET /api/v1/health
	
	// TODO: Implement additional health check endpoints in the future
	// health.Get("/live", healthHandler.CheckLiveness)      // GET /api/v1/health/live
	// health.Get("/ready", healthHandler.CheckReadiness)    // GET /api/v1/health/ready
	// health.Get("/detailed", healthHandler.CheckDetailed)  // GET /api/v1/health/detailed
} 
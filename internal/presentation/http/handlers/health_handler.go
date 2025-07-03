package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"trading-alchemist/internal/presentation/responses"
)

// HealthResponse represents the health check response
//
// swagger:model HealthResponse
type HealthResponse struct {
	// Service status
	// example: ok
	Status string `json:"status"`
	// Current timestamp
	// example: 2023-12-01T10:00:00Z
	Timestamp string `json:"timestamp"`
	// API version
	// example: 1.0.0
	Version string `json:"version"`
	// Service uptime (optional)
	// example: 24h30m
	Uptime string `json:"uptime,omitempty"`
}

// CheckHealth performs a health check
// @Summary Health check
// @Description Returns the current health status of the API. This endpoint can be used for monitoring and load balancer health checks.
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} responses.SuccessResponse{data=HealthResponse} "Service is healthy"
// @Failure 503 {object} responses.ErrorResponse "Service is unhealthy"
// @Router /health [get]
func CheckHealth(c *fiber.Ctx) error {
	// TODO: Add database connectivity check
	// TODO: Add external service health checks

	healthData := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
	}

	return responses.SendSuccess(c, healthData)
} 
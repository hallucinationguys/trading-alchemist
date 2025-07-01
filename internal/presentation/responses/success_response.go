package responses

import (
	"github.com/gofiber/fiber/v2"
)

// SuccessResponse represents a successful response
//
// swagger:model SuccessResponse
type SuccessResponse struct {
	// Response data
	Data interface{} `json:"data"`
	// Success indicator
	// example: true
	Success bool `json:"success"`
	// Optional success message
	// example: Operation completed successfully
	Message string `json:"message,omitempty"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(data interface{}, message ...string) *SuccessResponse {
	response := &SuccessResponse{
		Data:    data,
		Success: true,
	}

	if len(message) > 0 {
		response.Message = message[0]
	}

	return response
}

// SendSuccess sends a success response
func SendSuccess(c *fiber.Ctx, data interface{}, message ...string) error {
	response := NewSuccessResponse(data, message...)
	return c.JSON(response)
}

// SendCreated sends a 201 Created response
func SendCreated(c *fiber.Ctx, data interface{}, message ...string) error {
	response := NewSuccessResponse(data, message...)
	return c.Status(fiber.StatusCreated).JSON(response)
}

// SendAccepted sends a 202 Accepted response
func SendAccepted(c *fiber.Ctx, data interface{}, message ...string) error {
	response := NewSuccessResponse(data, message...)
	return c.Status(fiber.StatusAccepted).JSON(response)
}

// SendNoContent sends a 204 No Content response
func SendNoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
} 
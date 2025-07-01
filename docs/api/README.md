# Trading Alchemist API Documentation

This directory contains the API documentation for the Trading Alchemist project, generated using Swagger/OpenAPI 2.0.

## 🚀 Quick Start

### Viewing the API Documentation

1. **Start the API server:**
   ```bash
   # Using devbox (recommended)
   devbox run dev
   
   # Or using Docker
   make docker-dev
   
   # Or locally
   make run
   ```

2. **Access the documentation:**
   - **Swagger JSON:** `http://localhost:8080/swagger.json`
   - **Interactive UI:** `http://localhost:8080/docs` (redirects to Swagger UI)

### Regenerating Documentation

```bash
# Generate Swagger docs
devbox run swagger

# Or manually (requires swag tool installed)
make swagger-generate
```

## 📋 API Overview

The Trading Alchemist API provides secure authentication and user management capabilities:

### 🔐 Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/auth/send-magic-link` | Send magic link to user's email |
| `POST` | `/api/auth/verify` | Verify magic link token and get JWT |

### 👤 User Management Endpoints

| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| `GET` | `/api/users/profile` | Get current user profile | Required |
| `PUT` | `/api/users/profile` | Update user profile | Required |

### 🏥 Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | API health status |

## 🔑 Authentication

This API uses JWT Bearer token authentication:

```http
Authorization: Bearer <your-jwt-token>
```

### Getting a Token

1. Send a magic link to your email:
   ```bash
   curl -X POST http://localhost:8080/api/auth/send-magic-link \
     -H "Content-Type: application/json" \
     -d '{"email": "user@example.com"}'
   ```

2. Verify the magic link token:
   ```bash
   curl -X POST http://localhost:8080/api/auth/verify \
     -H "Content-Type: application/json" \
     -d '{"token": "your-magic-link-token"}'
   ```

## 📝 Request/Response Examples

### Send Magic Link

**Request:**
```json
{
  "email": "user@example.com",
  "purpose": "login"
}
```

**Response:**
```json
{
  "data": {
    "message": "If this email is registered, a magic link has been sent",
    "sent": true
  },
  "success": true,
  "message": "Magic link sent successfully"
}
```

### Verify Magic Link

**Request:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "data": {
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "user@example.com",
      "email_verified": true,
      "first_name": "John",
      "last_name": "Doe",
      "full_name": "John Doe",
      "display_name": "John Doe",
      "is_active": true,
      "created_at": "2023-12-01T10:00:00Z",
      "updated_at": "2023-12-01T10:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400
  },
  "success": true,
  "message": "Authentication successful"
}
```

### Get User Profile

**Request:**
```bash
curl -X GET http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer <your-jwt-token>"
```

**Response:**
```json
{
  "data": {
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "user@example.com",
      "email_verified": true,
      "first_name": "John",
      "last_name": "Doe",
      "avatar_url": "https://example.com/avatar.jpg",
      "full_name": "John Doe",
      "display_name": "John Doe",
      "is_active": true,
      "created_at": "2023-12-01T10:00:00Z",
      "updated_at": "2023-12-01T10:00:00Z"
    }
  },
  "success": true,
  "message": "Profile retrieved successfully"
}
```

## 🔧 Error Handling

All errors follow a consistent format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid email address",
    "details": "The email field must be a valid email address"
  },
  "success": false
}
```

### Common Error Codes

| Code | Description |
|------|-------------|
| `VALIDATION_ERROR` | Request validation failed |
| `UNAUTHORIZED` | Authentication required or invalid |
| `FORBIDDEN` | Access denied |
| `NOT_FOUND` | Resource not found |
| `CONFLICT` | Resource already exists |
| `INTERNAL_ERROR` | Server error |

## 🎯 Response Format

All API responses follow a consistent structure:

### Success Response
```json
{
  "data": {...},
  "success": true,
  "message": "Optional success message"
}
```

### Error Response
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message",
    "details": "Optional error details"
  },
  "success": false
}
```

## 🛠️ Development

### Adding New Endpoints

1. Create handler functions with Swagger annotations:
   ```go
   // @Summary Endpoint summary
   // @Description Detailed description
   // @Tags TagName
   // @Accept json
   // @Produce json
   // @Param request body dto.RequestType true "Request description"
   // @Success 200 {object} responses.SuccessResponse{data=dto.ResponseType}
   // @Failure 400 {object} responses.ErrorResponse
   // @Router /endpoint [method]
   func (h *Handler) HandlerFunction(c *fiber.Ctx) error {
       // Implementation
   }
   ```

2. Regenerate documentation:
   ```bash
   devbox run swagger
   ```

### Swagger Annotations Reference

- `@Summary` - Short endpoint description
- `@Description` - Detailed description
- `@Tags` - Group endpoints by functionality
- `@Accept` - Content types the endpoint accepts
- `@Produce` - Content types the endpoint produces
- `@Param` - Request parameters
- `@Success` - Success response definition
- `@Failure` - Error response definition
- `@Router` - Route path and HTTP method
- `@Security` - Authentication requirements

## 📚 Additional Resources

- [Swagger Documentation](https://swagger.io/docs/)
- [Go Swagger (swaggo)](https://github.com/swaggo/swag)
- [Fiber Framework](https://docs.gofiber.io/)
- [OpenAPI Specification](https://swagger.io/specification/)

## 🐛 Troubleshooting

### Swagger Generation Issues

1. **Empty swagger.json:** Ensure all handlers have proper annotations
2. **Missing endpoints:** Check that handlers are properly imported
3. **Build errors:** Verify Go module dependencies are up to date

### Common Solutions

```bash
# Clean and regenerate
make clean
devbox run install-tools
devbox run swagger

# Update dependencies
go mod tidy
go mod download
``` 
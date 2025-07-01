package dto

// SendMagicLinkRequest represents the request to send a magic link
//
// swagger:model SendMagicLinkRequest
type SendMagicLinkRequest struct {
	// Email address to send the magic link to
	// required: true
	// example: user@example.com
	Email string `json:"email" validate:"required,email"`
	// Purpose of the magic link (login, email_verification, password_reset)
	// example: login
	Purpose   string `json:"purpose,omitempty"` // login, email_verification, password_reset
	IPAddress string `json:"-"`                 // Set by middleware
	UserAgent string `json:"-"`                 // Set by middleware
}

// SendMagicLinkResponse represents the response after sending a magic link
//
// swagger:model SendMagicLinkResponse
type SendMagicLinkResponse struct {
	// Response message
	// example: If this email is registered, a magic link has been sent
	Message string `json:"message"`
	// Whether the email was sent successfully
	// example: true
	Sent bool `json:"sent"`
}

// VerifyMagicLinkRequest represents the request to verify a magic link
//
// swagger:model VerifyMagicLinkRequest
type VerifyMagicLinkRequest struct {
	// Magic link token to verify
	// required: true
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	Token string `json:"token" validate:"required"`
}

// VerifyMagicLinkResponse represents the response after verifying a magic link
//
// swagger:model VerifyMagicLinkResponse
type VerifyMagicLinkResponse struct {
	// User information
	User UserResponse `json:"user"`
	// JWT access token
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	AccessToken string `json:"access_token"`
	// Token type
	// example: Bearer
	TokenType string `json:"token_type"`
	// Token expiration in seconds
	// example: 86400
	ExpiresIn int64 `json:"expires_in"` // seconds
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Message string `json:"message"`
	Sent    bool   `json:"sent"`
}

// LogoutResponse represents a logout response
type LogoutResponse struct {
	Message string `json:"message"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse represents a refresh token response
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
} 
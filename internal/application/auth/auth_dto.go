package auth

// SendMagicLinkRequest represents the request to send a magic link
type SendMagicLinkRequest struct {
	// Email address to send the magic link to
	Email string `json:"email" validate:"required,email"`
	// Purpose of the magic link (login, email_verification, password_reset)
	Purpose   string `json:"purpose,omitempty"` // login, email_verification, password_reset
	IPAddress string `json:"-"`                 // Set by middleware
	UserAgent string `json:"-"`                 // Set by middleware
}

// SendMagicLinkResponse represents the response after sending a magic link
type SendMagicLinkResponse struct {
	// Response message
	Message string `json:"message"`
	// Whether the email was sent successfully
	Sent bool `json:"sent"`
}

// VerifyMagicLinkRequest represents the request to verify a magic link
type VerifyMagicLinkRequest struct {
	// Magic link token to verify
	Token string `json:"token" validate:"required"`
}

// VerifyMagicLinkResponse represents the response after verifying a magic link
type VerifyMagicLinkResponse struct {
	// User information
	User UserResponse `json:"user"`
	// JWT access token
	AccessToken string `json:"access_token"`
	// Token type
	TokenType string `json:"token_type"`
	// Token expiration in seconds
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
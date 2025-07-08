package auth

import (
	"time"

	"github.com/google/uuid"
)

type MagicLinkPurpose string

const (
	MagicLinkPurposeLogin             MagicLinkPurpose = "login"
	MagicLinkPurposeEmailVerification MagicLinkPurpose = "email_verification"
	MagicLinkPurposePasswordReset     MagicLinkPurpose = "password_reset"
)

type MagicLink struct {
	ID        uuid.UUID         `json:"id" db:"id"`
	UserID    uuid.UUID         `json:"user_id" db:"user_id"`
	Token     string            `json:"token" db:"token"`
	TokenHash string            `json:"-" db:"token_hash"` // Never expose in JSON
	ExpiresAt time.Time         `json:"expires_at" db:"expires_at"`
	UsedAt    *time.Time        `json:"used_at" db:"used_at"`
	IPAddress *string           `json:"ip_address" db:"ip_address"`
	UserAgent *string           `json:"user_agent" db:"user_agent"`
	Purpose   MagicLinkPurpose  `json:"purpose" db:"purpose"`
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
}

// IsExpired checks if the magic link has expired
func (ml *MagicLink) IsExpired() bool {
	return time.Now().After(ml.ExpiresAt)
}

// IsUsed checks if the magic link has already been used
func (ml *MagicLink) IsUsed() bool {
	return ml.UsedAt != nil
}

// IsValid checks if the magic link is valid (not expired and not used)
func (ml *MagicLink) IsValid() bool {
	return !ml.IsExpired() && !ml.IsUsed()
}

// MarkAsUsed marks the magic link as used
func (ml *MagicLink) MarkAsUsed() {
	now := time.Now()
	ml.UsedAt = &now
}

// GetPurposeString returns the string representation of the purpose
func (ml *MagicLink) GetPurposeString() string {
	return string(ml.Purpose)
} 
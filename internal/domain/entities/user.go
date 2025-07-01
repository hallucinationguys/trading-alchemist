package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	Email         string     `json:"email" db:"email"`
	EmailVerified bool       `json:"email_verified" db:"email_verified"`
	FirstName     *string    `json:"first_name" db:"first_name"`
	LastName      *string    `json:"last_name" db:"last_name"`
	AvatarURL     *string    `json:"avatar_url" db:"avatar_url"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// FullName returns the full name of the user
func (u *User) FullName() string {
	var fullName string
	if u.FirstName != nil {
		fullName = *u.FirstName
	}
	if u.LastName != nil {
		if fullName != "" {
			fullName += " "
		}
		fullName += *u.LastName
	}
	return fullName
}

// DisplayName returns the display name (full name or email)
func (u *User) DisplayName() string {
	if fullName := u.FullName(); fullName != "" {
		return fullName
	}
	return u.Email
}

// IsEmailVerified checks if the user's email is verified
func (u *User) IsEmailVerified() bool {
	return u.EmailVerified
}

// IsAccountActive checks if the user account is active
func (u *User) IsAccountActive() bool {
	return u.IsActive
} 
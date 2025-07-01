package utils

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail validates if the given string is a valid email address
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

// NormalizeEmail normalizes an email address by trimming whitespace and converting to lowercase
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// IsNotEmpty checks if a string is not empty after trimming whitespace
func IsNotEmpty(s string) bool {
	return strings.TrimSpace(s) != ""
}

// ValidateStringLength validates if a string length is within the specified range
func ValidateStringLength(s string, min, max int) bool {
	length := len(strings.TrimSpace(s))
	return length >= min && length <= max
} 
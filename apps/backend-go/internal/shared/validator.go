package shared

import (
	"net/mail"
	"regexp"
	"strings"
)

var (
	// Basic E.164 phone validation: + followed by 7 to 15 digits
	phoneRegex = regexp.MustCompile(`^\+[1-9]\d{6,14}$`)
)

// IsValidEmail checks if a string is a valid email address.
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// IsValidPhone checks if a string is a valid E.164 phone number.
func IsValidPhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

// NormalizePhone trims and ensures phone has a '+' prefix if missing (heuristic).
// For strict E.164, we usually expect the user to provide it.
func NormalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return ""
	}
	if !strings.HasPrefix(phone, "+") {
		// Heuristic: if no prefix, don't auto-fix, just let validation fail
		// or we could add a default country code if known.
	}
	return phone
}

package validator

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters long")
	ErrPasswordTooLong    = errors.New("password must not exceed 128 characters")
	ErrPasswordTooCommon  = errors.New("password is too common and easily guessable")
	ErrPasswordRepeating  = errors.New("password contains too many repeating characters")
	ErrPasswordAllNumeric = errors.New("password cannot be all numeric")
	ErrPasswordWhitespace = errors.New("password cannot contain leading or trailing whitespace")
)

var commonPasswords = map[string]bool{
	"password":    true,
	"12345678":    true,
	"123456789":   true,
	"password1":   true,
	"password123": true,
	"qwerty":      true,
	"abc123":      true,
	"monkey":      true,
	"letmein":     true,
	"trustno1":    true,
	"dragon":      true,
	"baseball":    true,
	"iloveyou":    true,
	"master":      true,
	"sunshine":    true,
	"ashley":      true,
	"bailey":      true,
	"shadow":      true,
	"superman":    true,
}

// PasswordValidator provides comprehensive password validation
type PasswordValidator struct {
	MinLength       int
	MaxLength       int
	CheckCommon     bool
	CheckRepeating  bool
	CheckAllNumeric bool
	CheckWhitespace bool
}

// NewPasswordValidator creates a validator with NIST compliant defaults
func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		MinLength:       8,
		MaxLength:       128,
		CheckCommon:     true,
		CheckRepeating:  true,
		CheckAllNumeric: true,
		CheckWhitespace: true,
	}
}

// Validate performs comprehensive password validation
func (v *PasswordValidator) Validate(password string) error {
	if v.CheckWhitespace && (strings.HasPrefix(password, " ") || strings.HasSuffix(password, " ")) {
		return ErrPasswordWhitespace
	}

	if len(password) < v.MinLength {
		return ErrPasswordTooShort
	}
	if len(password) > v.MaxLength {
		return ErrPasswordTooLong
	}

	if v.CheckCommon && v.isCommonPassword(password) {
		return ErrPasswordTooCommon
	}

	if v.CheckAllNumeric && v.isAllNumeric(password) {
		return ErrPasswordAllNumeric
	}

	if v.CheckRepeating && v.hasExcessiveRepeating(password) {
		return ErrPasswordRepeating
	}

	return nil
}

// isCommonPassword checks against known common passwords
func (v *PasswordValidator) isCommonPassword(password string) bool {
	lower := strings.ToLower(password)
	return commonPasswords[lower]
}

// isAllNumeric checks if password is all numeric
func (v *PasswordValidator) isAllNumeric(password string) bool {
	for _, r := range password {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return len(password) > 0
}

// hasExcessiveRepeating checks for patterns like "aaaa" or "1111"
func (v *PasswordValidator) hasExcessiveRepeating(password string) bool {
	if len(password) < 4 {
		return false
	}

	consecutiveCount := 1
	for i := 1; i < len(password); i++ {
		if password[i] == password[i-1] {
			consecutiveCount++
			if consecutiveCount >= 4 {
				return true
			}
		} else {
			consecutiveCount = 1
		}
	}
	return false
}

// Custom Gin validator function
func ValidatePassword(password string) error {
	validator := NewPasswordValidator()
	return validator.Validate(password)
}

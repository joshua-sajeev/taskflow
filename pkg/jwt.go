package pkg

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrEmptyUsername     = errors.New("username cannot be empty")
	ErrTokenCreation     = errors.New("failed to create token")
	ErrTokenValidation   = errors.New("token validation failed")
	ErrInvalidClaims     = errors.New("invalid token claims")
	ErrUsernameNotFound  = errors.New("username claim not found or invalid")
	ErrUnexpectedSigning = errors.New("unexpected signing method")
)

const (
	TokenExpiration    = 24 * time.Hour
	UsernameClaimKey   = "username"
	ExpirationClaimKey = "exp"
)

// CreateToken generates a JWT token for the given username
func CreateToken(username string, secretKey []byte) (string, error) {
	if username == "" {
		return "", ErrEmptyUsername
	}

	if len(secretKey) == 0 {
		return "", fmt.Errorf("%w: secret key cannot be empty", ErrTokenCreation)
	}

	claims := jwt.MapClaims{
		UsernameClaimKey:   username,
		ExpirationClaimKey: time.Now().Add(TokenExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrTokenCreation, err)
	}

	return tokenString, nil
}

// ValidateToken parses and validates a JWT token string
func ValidateToken(tokenString string, secretKey []byte) (*jwt.MapClaims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("%w: token string cannot be empty", ErrTokenValidation)
	}

	if len(secretKey) == 0 {
		return nil, fmt.Errorf("%w: secret key cannot be empty", ErrTokenValidation)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigning, token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTokenValidation, err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("%w: token is not valid", ErrTokenValidation)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("%w: cannot convert to MapClaims", ErrInvalidClaims)
	}

	return &claims, nil
}

// GetUsernameFromToken extracts the username from a JWT token
func GetUsernameFromToken(tokenString string, secretKey []byte) (string, error) {
	claims, err := ValidateToken(tokenString, secretKey)
	if err != nil {
		return "", err
	}

	username, ok := (*claims)[UsernameClaimKey].(string)
	if !ok {
		return "", ErrUsernameNotFound
	}

	if username == "" {
		return "", fmt.Errorf("%w: username is empty", ErrUsernameNotFound)
	}

	return username, nil
}

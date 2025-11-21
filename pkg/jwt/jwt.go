package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrEmptyUserID       = errors.New("user ID cannot be empty")
	ErrTokenCreation     = errors.New("failed to create token")
	ErrTokenValidation   = errors.New("token validation failed")
	ErrInvalidClaims     = errors.New("invalid token claims")
	ErrUserIDNotFound    = errors.New("user ID claim not found or invalid")
	ErrUnexpectedSigning = errors.New("unexpected signing method")
)

const (
	TokenExpiration    = 24 * time.Hour
	UserIDClaimKey     = "user_id"
	ExpirationClaimKey = "exp"
)

func CreateToken(userID int, email string, secretKey []byte) (string, error) {
	if userID == 0 {
		return "", ErrEmptyUserID
	}

	if len(secretKey) == 0 {
		return "", fmt.Errorf("%w: secret key cannot be empty", ErrTokenCreation)
	}

	claims := jwt.MapClaims{
		UserIDClaimKey:     userID,
		"email":            email,
		"iat":              time.Now().Unix(),
		ExpirationClaimKey: time.Now().Add(TokenExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrTokenCreation, err)
	}

	return tokenString, nil
}

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

func GetUserIDFromToken(tokenString string, secretKey []byte) (int, error) {
	claims, err := ValidateToken(tokenString, secretKey)
	if err != nil {
		return 0, err
	}

	userID, ok := (*claims)[UserIDClaimKey].(float64)
	if !ok {
		return 0, ErrUserIDNotFound
	}

	return int(userID), nil
}

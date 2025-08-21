package pkg

import (
	"fmt"
	"taskflow/internal/common"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(username string, secretKey []byte) (string, error) {
	if username == "" {
		return "", common.ErrorResponse{Message: "username cannot be null"}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(tokenString string, secretKey []byte) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &claims, nil
}

func GetUsernameFromToken(tokenString string, secretKey []byte) (string, error) {
	claims, err := ValidateToken(tokenString, secretKey)
	if err != nil {
		return "", err
	}

	username, ok := (*claims)["username"].(string)
	if !ok {
		return "", fmt.Errorf("username claim not found or invalid")
	}

	return username, nil
}

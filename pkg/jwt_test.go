package pkg

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	testSecretKey    = []byte("test-secret-key-123")
	invalidSecretKey = []byte("wrong-secret-key")
	emptySecretKey   = []byte("")
)

func TestCreateToken(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		secretKey     []byte
		expectedError error
		wantErr       bool
	}{
		{
			name:          "valid user",
			username:      "alice",
			secretKey:     testSecretKey,
			expectedError: nil,
			wantErr:       false,
		},
		{
			name:          "empty username",
			username:      "",
			secretKey:     testSecretKey,
			expectedError: ErrEmptyUsername,
			wantErr:       true,
		},
		{
			name:          "empty secret key",
			username:      "alice",
			secretKey:     emptySecretKey,
			expectedError: ErrTokenCreation,
			wantErr:       true,
		},
		{
			name:          "nil secret key",
			username:      "alice",
			secretKey:     nil,
			expectedError: ErrTokenCreation,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := CreateToken(tt.username, tt.secretKey)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateToken() expected error but got none")
					return
				}

				if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
					t.Errorf("CreateToken() error = %v, expected error type %v", err, tt.expectedError)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateToken() unexpected error = %v", err)
				return
			}

			if token == "" {
				t.Errorf("CreateToken() returned empty token")
				return
			}

			parsedToken, parseErr := jwt.Parse(token, func(token *jwt.Token) (any, error) {
				return tt.secretKey, nil
			})

			if parseErr != nil {
				t.Errorf("CreateToken() generated token could not be parsed: %v", parseErr)
				return
			}

			if !parsedToken.Valid {
				t.Errorf("CreateToken() generated invalid token")
				return
			}

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				t.Errorf("CreateToken() token claims are not MapClaims")
				return
			}

			if claims[UsernameClaimKey] != tt.username {
				t.Errorf("CreateToken() username in token = %v, want %v", claims[UsernameClaimKey], tt.username)
			}

			if exp, exists := claims[ExpirationClaimKey]; exists {
				if expFloat, ok := exp.(float64); ok {
					expTime := time.Unix(int64(expFloat), 0)
					if expTime.Before(time.Now()) {
						t.Errorf("CreateToken() token expiration is in the past")
					}
				} else {
					t.Errorf("CreateToken() expiration claim is not a valid timestamp")
				}
			} else {
				t.Errorf("CreateToken() token missing expiration claim")
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	validToken, err := CreateToken("johndoe", testSecretKey)
	if err != nil {
		t.Fatalf("failed to create valid token for testing: %v", err)
	}

	wrongSecretToken, err := CreateToken("johndoe", invalidSecretKey)
	if err != nil {
		t.Fatalf("failed to create token with wrong secret for testing: %v", err)
	}

	expiredClaims := jwt.MapClaims{
		UsernameClaimKey:   "johndoe",
		ExpirationClaimKey: time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
	}
	expiredTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredToken, err := expiredTokenObj.SignedString(testSecretKey)
	if err != nil {
		t.Fatalf("failed to create expired token for testing: %v", err)
	}

	tests := []struct {
		name          string
		tokenString   string
		secretKey     []byte
		expectedError error
		wantErr       bool
		wantUsername  string
	}{
		{
			name:          "valid token",
			tokenString:   validToken,
			secretKey:     testSecretKey,
			expectedError: nil,
			wantErr:       false,
			wantUsername:  "johndoe",
		},
		{
			name:          "invalid secret",
			tokenString:   wrongSecretToken,
			secretKey:     testSecretKey,
			expectedError: ErrTokenValidation,
			wantErr:       true,
		},
		{
			name:          "malformed token",
			tokenString:   "not-a-token",
			secretKey:     testSecretKey,
			expectedError: ErrTokenValidation,
			wantErr:       true,
		},
		{
			name:          "empty token string",
			tokenString:   "",
			secretKey:     testSecretKey,
			expectedError: ErrTokenValidation,
			wantErr:       true,
		},
		{
			name:          "empty secret key",
			tokenString:   validToken,
			secretKey:     emptySecretKey,
			expectedError: ErrTokenValidation,
			wantErr:       true,
		},
		{
			name:          "expired token",
			tokenString:   expiredToken,
			secretKey:     testSecretKey,
			expectedError: ErrTokenValidation,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.tokenString, tt.secretKey)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateToken() expected error but got none")
					return
				}

				if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
					t.Errorf("ValidateToken() error = %v, expected error type %v", err, tt.expectedError)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateToken() unexpected error = %v", err)
				return
			}

			if claims == nil {
				t.Errorf("ValidateToken() returned nil claims")
				return
			}

			if tt.wantUsername != "" {
				if username := (*claims)[UsernameClaimKey]; username != tt.wantUsername {
					t.Errorf("ValidateToken() username = %v, want %v", username, tt.wantUsername)
				}
			}
		})
	}
}

func TestGetUsernameFromToken(t *testing.T) {
	validToken, err := CreateToken("alice", testSecretKey)
	if err != nil {
		t.Fatalf("failed to create valid token for testing: %v", err)
	}

	claimsWithoutUsername := jwt.MapClaims{
		ExpirationClaimKey: time.Now().Add(time.Hour).Unix(),
		"other_claim":      "some_value",
	}
	tokenWithoutUsernameObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsWithoutUsername)
	tokenWithoutUsername, err := tokenWithoutUsernameObj.SignedString(testSecretKey)
	if err != nil {
		t.Fatalf("failed to create token without username for testing: %v", err)
	}

	claimsWithEmptyUsername := jwt.MapClaims{
		UsernameClaimKey:   "",
		ExpirationClaimKey: time.Now().Add(time.Hour).Unix(),
	}
	tokenWithEmptyUsernameObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsWithEmptyUsername)
	tokenWithEmptyUsername, err := tokenWithEmptyUsernameObj.SignedString(testSecretKey)
	if err != nil {
		t.Fatalf("failed to create token with empty username for testing: %v", err)
	}

	tests := []struct {
		name          string
		tokenString   string
		secretKey     []byte
		expectedError error
		wantErr       bool
		wantUsername  string
	}{
		{
			name:          "valid token",
			tokenString:   validToken,
			secretKey:     testSecretKey,
			expectedError: nil,
			wantErr:       false,
			wantUsername:  "alice",
		},
		{
			name:          "invalid token string",
			tokenString:   "this.is.not.jwt",
			secretKey:     testSecretKey,
			expectedError: ErrTokenValidation,
			wantErr:       true,
		},
		{
			name:          "token without username claim",
			tokenString:   tokenWithoutUsername,
			secretKey:     testSecretKey,
			expectedError: ErrUsernameNotFound,
			wantErr:       true,
		},
		{
			name:          "token with empty username",
			tokenString:   tokenWithEmptyUsername,
			secretKey:     testSecretKey,
			expectedError: ErrUsernameNotFound,
			wantErr:       true,
		},
		{
			name:          "empty token string",
			tokenString:   "",
			secretKey:     testSecretKey,
			expectedError: ErrTokenValidation,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			username, err := GetUsernameFromToken(tt.tokenString, tt.secretKey)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetUsernameFromToken() expected error but got none")
					return
				}

				if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
					t.Errorf("GetUsernameFromToken() error = %v, expected error type %v", err, tt.expectedError)
				}
				return
			}

			if err != nil {
				t.Errorf("GetUsernameFromToken() unexpected error = %v", err)
				return
			}

			if username != tt.wantUsername {
				t.Errorf("GetUsernameFromToken() = %v, want %v", username, tt.wantUsername)
			}
		})
	}
}

func BenchmarkCreateToken(b *testing.B) {
	username := "testuser"
	secretKey := testSecretKey

	for b.Loop() {
		_, err := CreateToken(username, secretKey)
		if err != nil {
			b.Fatalf("CreateToken failed: %v", err)
		}
	}
}

func BenchmarkValidateToken(b *testing.B) {
	token, err := CreateToken("testuser", testSecretKey)
	if err != nil {
		b.Fatalf("failed to create token for benchmark: %v", err)
	}

	for b.Loop() {
		_, err := ValidateToken(token, testSecretKey)
		if err != nil {
			b.Fatalf("ValidateToken failed: %v", err)
		}
	}
}

func BenchmarkGetUsernameFromToken(b *testing.B) {
	token, err := CreateToken("testuser", testSecretKey)
	if err != nil {
		b.Fatalf("failed to create token for benchmark: %v", err)
	}

	for b.Loop() {
		_, err := GetUsernameFromToken(token, testSecretKey)
		if err != nil {
			b.Fatalf("GetUsernameFromToken failed: %v", err)
		}
	}
}

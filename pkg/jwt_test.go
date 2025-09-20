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
		userID        int
		email         string
		secretKey     []byte
		expectedError error
		wantErr       bool
	}{
		{
			name:      "valid token",
			userID:    1,
			email:     "alice@example.com",
			secretKey: testSecretKey,
			wantErr:   false,
		},
		{
			name:          "empty userID",
			userID:        0,
			email:         "alice@example.com",
			secretKey:     testSecretKey,
			expectedError: ErrEmptyUserID,
			wantErr:       true,
		},
		{
			name:          "empty secret key",
			userID:        1,
			email:         "bob@example.com",
			secretKey:     emptySecretKey,
			expectedError: ErrTokenCreation,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := CreateToken(tt.userID, tt.email, tt.secretKey)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateToken() expected error, got nil")
					return
				}
				if !errors.Is(err, tt.expectedError) {
					t.Errorf("CreateToken() error = %v, want %v", err, tt.expectedError)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateToken() unexpected error: %v", err)
				return
			}
			if token == "" {
				t.Errorf("CreateToken() returned empty token")
			}

			// Parse back to verify claims
			parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
				return tt.secretKey, nil
			})
			if err != nil || !parsed.Valid {
				t.Errorf("CreateToken() produced invalid token: %v", err)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	validToken, _ := CreateToken(42, "john@example.com", testSecretKey)

	expiredClaims := jwt.MapClaims{
		UserIDClaimKey:     42,
		"email":            "john@example.com",
		ExpirationClaimKey: time.Now().Add(-time.Hour).Unix(),
	}
	expiredTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredToken, _ := expiredTokenObj.SignedString(testSecretKey)

	tests := []struct {
		name       string
		tokenStr   string
		secretKey  []byte
		wantErr    bool
		wantUserID int
	}{
		{
			name:       "valid token",
			tokenStr:   validToken,
			secretKey:  testSecretKey,
			wantErr:    false,
			wantUserID: 42,
		},
		{
			name:      "wrong secret",
			tokenStr:  validToken,
			secretKey: invalidSecretKey,
			wantErr:   true,
		},
		{
			name:      "malformed token",
			tokenStr:  "not-a-token",
			secretKey: testSecretKey,
			wantErr:   true,
		},
		{
			name:      "expired token",
			tokenStr:  expiredToken,
			secretKey: testSecretKey,
			wantErr:   true,
		},
		{
			name:      "empty token",
			tokenStr:  "",
			secretKey: testSecretKey,
			wantErr:   true,
		},
		{
			name:      "empty secret key",
			tokenStr:  validToken,
			secretKey: emptySecretKey,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.tokenStr, tt.secretKey)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateToken() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("ValidateToken() unexpected error: %v", err)
				return
			}
			if uid := int((*claims)[UserIDClaimKey].(float64)); uid != tt.wantUserID {
				t.Errorf("ValidateToken() userID = %v, want %v", uid, tt.wantUserID)
			}
		})
	}
}

func TestGetUserIDFromToken(t *testing.T) {
	validToken, _ := CreateToken(7, "alice@example.com", testSecretKey)

	noUserIDClaims := jwt.MapClaims{
		"email":            "alice@example.com",
		ExpirationClaimKey: time.Now().Add(time.Hour).Unix(),
	}
	noUserIDToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, noUserIDClaims).SignedString(testSecretKey)

	tests := []struct {
		name       string
		tokenStr   string
		secretKey  []byte
		wantErr    bool
		wantUserID int
	}{
		{
			name:       "valid token",
			tokenStr:   validToken,
			secretKey:  testSecretKey,
			wantErr:    false,
			wantUserID: 7,
		},
		{
			name:      "token without userID",
			tokenStr:  noUserIDToken,
			secretKey: testSecretKey,
			wantErr:   true,
		},
		{
			name:      "invalid token string",
			tokenStr:  "invalid.token.here",
			secretKey: testSecretKey,
			wantErr:   true,
		},
		{
			name:      "empty token string",
			tokenStr:  "",
			secretKey: testSecretKey,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uid, err := GetUserIDFromToken(tt.tokenStr, tt.secretKey)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetUserIDFromToken() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("GetUserIDFromToken() unexpected error: %v", err)
				return
			}
			if uid != tt.wantUserID {
				t.Errorf("GetUserIDFromToken() = %v, want %v", uid, tt.wantUserID)
			}
		})
	}
}

func BenchmarkCreateToken(b *testing.B) {
	for b.Loop() {
		if _, err := CreateToken(1, "bench@example.com", testSecretKey); err != nil {
			b.Fatalf("CreateToken failed: %v", err)
		}
	}
}

func BenchmarkValidateToken(b *testing.B) {
	token, _ := CreateToken(1, "bench@example.com", testSecretKey)
	for b.Loop() {
		if _, err := ValidateToken(token, testSecretKey); err != nil {
			b.Fatalf("ValidateToken failed: %v", err)
		}
	}
}

func BenchmarkGetUserIDFromToken(b *testing.B) {
	token, _ := CreateToken(1, "bench@example.com", testSecretKey)
	for b.Loop() {
		if _, err := GetUserIDFromToken(token, testSecretKey); err != nil {
			b.Fatalf("GetUserIDFromToken failed: %v", err)
		}
	}
}

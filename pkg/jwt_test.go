package pkg

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("test-secret")
var invalidSecret = []byte("wrongsecret")

func TestCreateToken(t *testing.T) {

	tests := []struct {
		name      string
		username  string
		secretKey []byte
		want      string
		wantErr   bool
	}{
		{"valid user", "alice", secretKey, "alice", false},
		{"empty username", "", secretKey, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := CreateToken(tt.username, tt.secretKey)

			if (gotErr != nil) != tt.wantErr {
				t.Fatalf("CreateToken(%q) error = %v, wantErr %v", tt.username, gotErr, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if got == "" {
				t.Errorf("CreateToken(%q) returned empty token", tt.username)
				return
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "johndoe",
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	validTokenStr, err := validToken.SignedString(secretKey)
	if err != nil {
		t.Fatalf("failed to create valid token: %v", err)
	}

	invalidToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "johndoe",
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	invalidTokenStr, err := invalidToken.SignedString(invalidSecret)
	if err != nil {
		t.Fatalf("failed to create invalid token: %v", err)
	}

	tests := []struct {
		name        string
		tokenString string
		secretKey   []byte
		want        jwt.MapClaims
		wantErr     bool
	}{
		{
			name:        "valid token",
			tokenString: validTokenStr,
			secretKey:   secretKey,
			want:        jwt.MapClaims{"username": "johndoe"},
			wantErr:     false,
		},
		{
			name:        "invalid secret",
			tokenString: invalidTokenStr,
			secretKey:   secretKey,
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "malformed token",
			tokenString: "not-a-token",
			secretKey:   secretKey,
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := ValidateToken(tt.tokenString, tt.secretKey)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ValidateToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ValidateToken() succeeded unexpectedly")
			}

			if tt.want != nil {
				if (*got)["username"] != tt.want["username"] {
					t.Errorf("ValidateToken() username = %v, want %v", (*got)["username"], tt.want["username"])
				}
			}
		})
	}
}

func TestGetUsernameFromToken(t *testing.T) {

	validToken, _ := CreateToken("alice", secretKey)
	tests := []struct {
		name        string
		tokenString string
		secretKey   []byte
		want        string
		wantErr     bool
	}{
		{"valid token", validToken, secretKey, "alice", false},
		{"invalid token string", "this.is.not.jwt", invalidSecret, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := GetUsernameFromToken(tt.tokenString, tt.secretKey)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetUsernameFromToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetUsernameFromToken() succeeded unexpectedly")
			}
			if got != "alice" {
				t.Errorf("GetUsernameFromToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

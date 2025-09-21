package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"taskflow/internal/common"
	"taskflow/pkg"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func setupGinTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name             string
		authHeader       string
		expectedStatus   int
		expectedResponse any
		shouldCallNext   bool
	}{
		// Note: For a real success case, you'd need a valid JWT token
		// This test focuses on the error cases that we can reliably test
		{
			name:           "failure case - missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: common.ErrorResponse{
				Message: "authorization header required",
			},
			shouldCallNext: false,
		},
		{
			name:           "failure case - invalid header format (no Bearer)",
			authHeader:     "invalid-token-123",
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: common.ErrorResponse{
				Message: "invalid authorization header format",
			},
			shouldCallNext: false,
		},
		{
			name:           "failure case - invalid header format (missing token)",
			authHeader:     "Bearer",
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: common.ErrorResponse{
				Message: "invalid authorization header format",
			},
			shouldCallNext: false,
		},
		{
			name:           "failure case - only Bearer with space",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: common.ErrorResponse{
				Message: "invalid or expired token",
			},
			shouldCallNext: false,
		},
		{
			name:           "failure case - wrong auth type",
			authHeader:     "Basic dGVzdDp0ZXN0",
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: common.ErrorResponse{
				Message: "invalid authorization header format",
			},
			shouldCallNext: false,
		},
		{
			name:           "failure case - multiple spaces",
			authHeader:     "Bearer  token-with-spaces",
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: common.ErrorResponse{
				Message: "invalid or expired token",
			},
			shouldCallNext: false,
		},
		{
			name:           "failure case - invalid token (will fail validation)",
			authHeader:     "Bearer invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: common.ErrorResponse{
				Message: "invalid or expired token",
			},
			shouldCallNext: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup router
			router := setupGinTest()

			var nextCalled bool

			userAuth := NewUserAuth("test-secret")
			router.Use(userAuth.AuthMiddleware())
			router.GET("/test", func(c *gin.Context) {
				nextCalled = true
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.shouldCallNext, nextCalled)

			if tt.expectedResponse != nil {
				var responseBody any
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				expectedBodyBytes, err := json.Marshal(tt.expectedResponse)
				assert.NoError(t, err)
				var expectedResponseBody any
				err = json.Unmarshal(expectedBodyBytes, &expectedResponseBody)
				assert.NoError(t, err)

				assert.Equal(t, expectedResponseBody, responseBody)
			}
		})
	}
}

func TestOptionalAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		shouldCallNext bool
		description    string
	}{
		// Note: For testing valid tokens, you'd need actual valid JWTs from your system
		{
			name:           "success case - no authorization header (continues)",
			authHeader:     "",
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
			description:    "Missing auth header should continue without setting userID",
		},
		{
			name:           "success case - invalid header format (continues)",
			authHeader:     "invalid-token-123",
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
			description:    "Invalid format should continue without setting userID",
		},
		{
			name:           "success case - wrong auth type (continues)",
			authHeader:     "Basic dGVzdDp0ZXN0",
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
			description:    "Wrong auth type should continue without setting userID",
		},
		{
			name:           "success case - Bearer only (continues)",
			authHeader:     "Bearer",
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
			description:    "Bearer without token should continue without setting userID",
		},
		{
			name:           "success case - invalid token (continues)",
			authHeader:     "Bearer invalid.token.here",
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
			description:    "Invalid token should continue without setting userID",
		},
		{
			name:           "success case - Bearer with space only",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
			description:    "Bearer with empty token should continue without setting userID",
		},
		{
			name:           "success case - multiple spaces",
			authHeader:     "Bearer  token-with-spaces",
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
			description:    "Multiple spaces should continue (token validation will fail gracefully)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup router
			router := setupGinTest()

			var nextCalled bool

			userAuth := NewUserAuth("test-secret")
			router.Use(userAuth.OptionalAuthMiddleware())
			router.GET("/test", func(c *gin.Context) {
				nextCalled = true
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code, tt.description)
			assert.Equal(t, tt.shouldCallNext, nextCalled, tt.description)

			// Optional auth middleware should always return 200 and call next
			assert.Equal(t, http.StatusOK, w.Code)
			assert.True(t, nextCalled)
		})
	}
}

// Integration-style test that tests the actual middleware behavior
func TestAuthMiddleware_Integration(t *testing.T) {
	tests := []struct {
		name         string
		setupRouter  func() *gin.Engine
		authHeader   string
		expectedCode int
		expectUserID bool
	}{
		{
			name: "auth middleware blocks invalid requests",
			setupRouter: func() *gin.Engine {
				r := setupGinTest()
				userAuth := NewUserAuth("test-secret")
				r.Use(userAuth.AuthMiddleware())
				r.GET("/protected", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "protected resource"})
				})
				return r
			},
			authHeader:   "",
			expectedCode: http.StatusUnauthorized,
			expectUserID: false,
		},
		{
			name: "optional auth middleware allows all requests",
			setupRouter: func() *gin.Engine {
				r := setupGinTest()
				userAuth := NewUserAuth("test-secret")
				r.Use(userAuth.OptionalAuthMiddleware())
				r.GET("/public", func(c *gin.Context) {
					userID, exists := c.Get("userID")
					response := gin.H{"message": "public resource"}
					if exists {
						response["userID"] = userID
					}
					c.JSON(http.StatusOK, response)
				})
				return r
			},
			authHeader:   "",
			expectedCode: http.StatusOK,
			expectUserID: false,
		},
		{
			name: "middleware handles malformed Bearer token",
			setupRouter: func() *gin.Engine {
				r := setupGinTest()
				userAuth := NewUserAuth("test-secret")
				r.Use(userAuth.AuthMiddleware())
				r.GET("/test", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "success"})
				})
				return r
			},
			authHeader:   "Bearer",
			expectedCode: http.StatusUnauthorized,
			expectUserID: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := tt.setupRouter()

			var endpoint string
			switch tt.name {
			case "auth middleware blocks invalid requests":
				endpoint = "/protected"
			case "optional auth middleware allows all requests":
				endpoint = "/public"
			default:
				endpoint = "/test"
			}

			req := httptest.NewRequest(http.MethodGet, endpoint, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
func TestAuthMiddleware_EdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secretKey := []byte("secret-key")
	validToken, _ := pkg.CreateToken(123, "test@example.com", secretKey)

	// Token with invalid signing method
	wrongSignToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user_id": 1,
	})
	wrongSignTokenString, _ := wrongSignToken.SignedString([]byte("secret-key"))

	// Token with non-numeric user_id
	userIDStrToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "not-a-number",
	})
	userIDStrTokenString, _ := userIDStrToken.SignedString(secretKey)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectNext     bool
	}{
		{
			name:           "valid token calls next",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectNext:     true,
		},
		{
			name:           "invalid signing method",
			authHeader:     "Bearer " + wrongSignTokenString,
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name:           "user_id not a number",
			authHeader:     "Bearer " + userIDStrTokenString,
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
		{
			name:           "empty token string",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupGinTest()
			var nextCalled bool
			userAuth := NewUserAuth("secret-key")
			router.Use(userAuth.AuthMiddleware())
			router.GET("/test", func(c *gin.Context) {
				nextCalled = true
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectNext, nextCalled)
		})
	}
}

func TestOptionalAuthMiddleware_EdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secretKey := []byte("secret-key")
	validToken, _ := pkg.CreateToken(456, "test@example.com", secretKey)

	tests := []struct {
		name       string
		authHeader string
		expectNext bool
	}{
		{"no auth header", "", true},
		{"invalid format", "invalid-token", true},
		{"wrong auth type", "Basic abc123", true},
		{"Bearer only", "Bearer", true},
		{"invalid token", "Bearer invalid.token", true},
		{"valid token", "Bearer " + validToken, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupGinTest()
			var nextCalled bool
			userAuth := NewUserAuth("test-secret")
			router.Use(userAuth.OptionalAuthMiddleware())
			router.GET("/test", func(c *gin.Context) {
				nextCalled = true
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, true, nextCalled)
		})
	}
}

package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"taskflow/internal/common"
	"taskflow/internal/domain/user"
	"taskflow/internal/repository/gorm/gorm_user"
	"taskflow/pkg/jwt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupGinTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestEnhancedAuthMiddleware_UserValidation(t *testing.T) {
	secretKey := []byte("test-secret")

	// Create a valid token for user ID 1
	validToken, err := jwt.CreateToken(1, "test@example.com", secretKey)
	assert.NoError(t, err)

	// Create a valid token for user ID 999 (non-existent user)
	deletedUserToken, err := jwt.CreateToken(999, "deleted@example.com", secretKey)
	assert.NoError(t, err)

	tests := []struct {
		name             string
		authHeader       string
		setupMock        func(*gorm_user.MockUserRepository)
		expectedStatus   int
		expectedResponse interface{}
		shouldCallNext   bool
	}{
		{
			name:       "success - user exists",
			authHeader: "Bearer " + validToken,
			setupMock: func(mockRepo *gorm_user.MockUserRepository) {
				mockRepo.On("GetByID", 1).Return(&user.User{
					ID:    1,
					Email: "test@example.com",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:       "failure - user deleted/not found",
			authHeader: "Bearer " + deletedUserToken,
			setupMock: func(mockRepo *gorm_user.MockUserRepository) {
				mockRepo.On("GetByID", 999).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: common.ErrorResponse{
				Message: "user account not found",
			},
			shouldCallNext: false,
		},
		{
			name:       "failure - database error during user lookup",
			authHeader: "Bearer " + validToken,
			setupMock: func(mockRepo *gorm_user.MockUserRepository) {
				mockRepo.On("GetByID", 1).Return(nil, errors.New("database connection error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: common.ErrorResponse{
				Message: "authentication failed",
			},
			shouldCallNext: false,
		},
		{
			name:           "failure - missing authorization header",
			authHeader:     "",
			setupMock:      func(mockRepo *gorm_user.MockUserRepository) {},
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: common.ErrorResponse{
				Message: "authorization header required",
			},
			shouldCallNext: false,
		},
		{
			name:           "failure - invalid token format",
			authHeader:     "Bearer invalid.token.here",
			setupMock:      func(mockRepo *gorm_user.MockUserRepository) {},
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: common.ErrorResponse{
				Message: "invalid or expired token",
			},
			shouldCallNext: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock repository
			mockRepo := new(gorm_user.MockUserRepository)
			tt.setupMock(mockRepo)

			// Setup router
			router := setupGinTest()
			var nextCalled bool

			userAuth := NewUserAuth("test-secret", mockRepo)
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
				var responseBody interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				expectedBodyBytes, err := json.Marshal(tt.expectedResponse)
				assert.NoError(t, err)
				var expectedResponseBody interface{}
				err = json.Unmarshal(expectedBodyBytes, &expectedResponseBody)
				assert.NoError(t, err)

				assert.Equal(t, expectedResponseBody, responseBody)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEnhancedOptionalAuthMiddleware_UserValidation(t *testing.T) {
	secretKey := []byte("test-secret")

	// Create tokens
	validToken, err := jwt.CreateToken(1, "test@example.com", secretKey)
	assert.NoError(t, err)

	deletedUserToken, err := jwt.CreateToken(999, "deleted@example.com", secretKey)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		setupMock      func(*gorm_user.MockUserRepository)
		expectedStatus int
		shouldCallNext bool
		shouldSetUser  bool
	}{
		{
			name:       "success - valid user",
			authHeader: "Bearer " + validToken,
			setupMock: func(mockRepo *gorm_user.MockUserRepository) {
				mockRepo.On("GetByID", 1).Return(&user.User{
					ID:    1,
					Email: "test@example.com",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
			shouldSetUser:  true,
		},
		{
			name:       "success - deleted user (continues without setting userID)",
			authHeader: "Bearer " + deletedUserToken,
			setupMock: func(mockRepo *gorm_user.MockUserRepository) {
				mockRepo.On("GetByID", 999).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
			shouldSetUser:  false,
		},
		{
			name:           "success - no auth header",
			authHeader:     "",
			setupMock:      func(mockRepo *gorm_user.MockUserRepository) {},
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
			shouldSetUser:  false,
		},
		{
			name:           "success - invalid token (continues)",
			authHeader:     "Bearer invalid.token",
			setupMock:      func(mockRepo *gorm_user.MockUserRepository) {},
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
			shouldSetUser:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock repository
			mockRepo := new(gorm_user.MockUserRepository)
			tt.setupMock(mockRepo)

			// Setup router
			router := setupGinTest()
			var nextCalled bool
			var userIDSet bool

			userAuth := NewUserAuth("test-secret", mockRepo)
			router.Use(userAuth.OptionalAuthMiddleware())
			router.GET("/test", func(c *gin.Context) {
				nextCalled = true
				_, exists := c.Get("userID")
				userIDSet = exists
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
			assert.Equal(t, tt.shouldSetUser, userIDSet)

			mockRepo.AssertExpectations(t)
		})
	}
}

// Integration test demonstrating the security fix
func TestSecurityScenario_DeletedUserCannotAccess(t *testing.T) {
	secretKey := []byte("test-secret")

	// Create a token for user ID 1
	userToken, err := jwt.CreateToken(1, "user@example.com", secretKey)
	assert.NoError(t, err)

	// Setup mock repository to simulate deleted user
	mockRepo := new(gorm_user.MockUserRepository)
	mockRepo.On("GetByID", 1).Return(nil, gorm.ErrRecordNotFound)

	// Setup router with protected endpoint
	router := setupGinTest()
	userAuth := NewUserAuth("test-secret", mockRepo)

	protected := router.Group("/api")
	protected.Use(userAuth.AuthMiddleware())
	{
		protected.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "This should not be accessible!"})
		})
	}

	// Try to access protected route with token from deleted user
	req := httptest.NewRequest(http.MethodGet, "/api/protected", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should be denied
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response common.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user account not found", response.Message)

	mockRepo.AssertExpectations(t)
}

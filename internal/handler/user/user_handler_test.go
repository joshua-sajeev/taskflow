package user_handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"taskflow/internal/common"
	"taskflow/internal/dto"
	user_service "taskflow/internal/service/user"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Setup Gin for tests
func setupGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestUserHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    any
		setupMock      func() *user_service.UserServiceMock
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "success case",
			requestBody: dto.CreateUserRequest{
				Email:    "test@example.com",
				Password: "password",
			},
			setupMock: func() *user_service.UserServiceMock {
				mockSvc := new(user_service.UserServiceMock)
				mockSvc.On("CreateUser", mock.Anything).Return(&dto.CreateUserResponse{
					ID:    1,
					Email: "test@example.com",
				}, nil)
				return mockSvc
			},
			expectedStatus: http.StatusCreated,
			expectedBody: dto.CreateUserResponse{
				ID:    1,
				Email: "test@example.com",
			},
		},
		{
			name: "failure case - email exists",
			requestBody: dto.CreateUserRequest{
				Email:    "exist@example.com",
				Password: "password",
			},
			setupMock: func() *user_service.UserServiceMock {
				mockSvc := new(user_service.UserServiceMock)
				mockSvc.On("CreateUser", mock.Anything).Return(nil, errors.New("email already exists"))
				return mockSvc
			},
			expectedStatus: http.StatusConflict,
			expectedBody: common.ErrorResponse{
				Message: "email already exists",
			},
		},
		{
			name:           "failure case - invalid JSON",
			requestBody:    `{"email":}`, // malformed
			setupMock:      func() *user_service.UserServiceMock { return new(user_service.UserServiceMock) },
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "invalid character '}' looking for beginning of value",
			},
		},
		{
			name: "failure case - No Email",
			requestBody: dto.CreateUserRequest{
				Password: "password",
			},
			setupMock:      func() *user_service.UserServiceMock { return new(user_service.UserServiceMock) },
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag",
			},
		},
		{
			name: "failure case - No Password",
			requestBody: dto.CreateUserRequest{
				Email: "test@gmail.com",
			},
			setupMock:      func() *user_service.UserServiceMock { return new(user_service.UserServiceMock) },
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Key: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			},
		},
		{
			name: "failure case - weak password;gin failure",
			requestBody: dto.CreateUserRequest{
				Email:    "weak@test.com",
				Password: "123",
			},
			setupMock: func() *user_service.UserServiceMock {
				mockSvc := new(user_service.UserServiceMock)
				return mockSvc
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Key: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag",
			},
		},
		{
			name: "failure case - weak password",
			requestBody: dto.CreateUserRequest{
				Email:    "weak@test.com",
				Password: "12345678",
			},
			setupMock: func() *user_service.UserServiceMock {
				mockSvc := new(user_service.UserServiceMock)

				mockSvc.On("CreateUser", mock.Anything).Return(nil, errors.New("password validation failed, choose a stronger password"))
				return mockSvc
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "password validation failed, choose a stronger password",
			},
		},
		{
			name: "failure case - hashing error",
			requestBody: dto.CreateUserRequest{
				Email:    "hash@test.com",
				Password: "StrongPass1!",
			},
			setupMock: func() *user_service.UserServiceMock {
				mockSvc := new(user_service.UserServiceMock)
				mockSvc.On("CreateUser", mock.Anything).Return(nil, errors.New("failed to hash password: internal error"))
				return mockSvc
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "failed to hash password: internal error",
			},
		},
		{
			name: "failure case - repo create generic error",
			requestBody: dto.CreateUserRequest{
				Email:    "repo@test.com",
				Password: "StrongPass1!",
			},
			setupMock: func() *user_service.UserServiceMock {
				mockSvc := new(user_service.UserServiceMock)
				mockSvc.On("CreateUser", mock.Anything).Return(nil, errors.New("failed to create user: db down"))
				return mockSvc
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "failed to create user: db down",
			},
		},
		{
			name:        "failure case - empty request body",
			requestBody: "",
			setupMock: func() *user_service.UserServiceMock {
				return new(user_service.UserServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Request body cannot be empty",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := tt.setupMock()
			handler := NewUserHandler(mockSvc, nil)

			router := setupGin()
			router.POST("/auth/register", handler.Register)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody any
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedBytes, _ := json.Marshal(tt.expectedBody)
			var expectedResponse any
			_ = json.Unmarshal(expectedBytes, &expectedResponse)

			assert.Equal(t, expectedResponse, responseBody)
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {

	tests := []struct {
		name           string
		requestBody    any
		setupMock      func() *user_service.UserServiceMock
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "success case",
			requestBody: dto.AuthRequest{
				Email:    "test@example.com",
				Password: "password",
			},
			setupMock: func() *user_service.UserServiceMock {
				mockSvc := new(user_service.UserServiceMock)
				mockSvc.
					On("AuthenticateUser", mock.Anything).
					Return(&dto.AuthResponse{
						ID:    1,
						Email: "test@example.com",
						Token: "jwt_token_here",
					}, nil)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			expectedBody: dto.AuthResponse{
				ID:    1,
				Email: "test@example.com",
				Token: "jwt_token_here",
			},
		},

		{
			name:        "failure - empty request body",
			requestBody: ``,
			setupMock: func() *user_service.UserServiceMock {
				return new(user_service.UserServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Request body cannot be empty",
			},
		},

		{
			name:        "failure - invalid json",
			requestBody: `{"email":}`,
			setupMock: func() *user_service.UserServiceMock {
				return new(user_service.UserServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "invalid character '}' looking for beginning of value",
			},
		},

		{
			name: "failure - missing email",
			requestBody: dto.AuthRequest{
				Password: "password",
			},
			setupMock: func() *user_service.UserServiceMock {
				return new(user_service.UserServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Key: 'AuthRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag",
			},
		},

		{
			name: "failure - missing password",
			requestBody: dto.AuthRequest{
				Email: "test@example.com",
			},
			setupMock: func() *user_service.UserServiceMock {
				return new(user_service.UserServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Key: 'AuthRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			},
		},

		{
			name: "failure - invalid credentials",
			requestBody: dto.AuthRequest{
				Email:    "wrong@example.com",
				Password: "wrongpw",
			},
			setupMock: func() *user_service.UserServiceMock {
				mockSvc := new(user_service.UserServiceMock)
				mockSvc.On("AuthenticateUser", mock.Anything).
					Return(nil, errors.New("invalid credentials"))
				return mockSvc
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: common.ErrorResponse{
				Message: "invalid credentials",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockService := tt.setupMock()

			handler := NewUserHandler(mockService, nil)
			router := setupGin()
			router.POST("/auth/login", handler.Login)

			var body []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var actual any
			err = json.Unmarshal(w.Body.Bytes(), &actual)
			assert.NoError(t, err)

			expectedBytes, _ := json.Marshal(tt.expectedBody)
			var expected any
			_ = json.Unmarshal(expectedBytes, &expected)

			assert.Equal(t, expected, actual)

			mockService.AssertExpectations(t)
		})
	}

}

func TestUserHandler_UpdatePassword(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    any
		setupMock      func() *user_service.UserServiceMock
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "success case",
			requestBody: dto.UpdatePasswordRequest{
				ID:          1,
				OldPassword: "superoldpassword",
				NewPassword: "ReallyStrongPass123!",
			},
			setupMock: func() *user_service.UserServiceMock {
				mockSvc := new(user_service.UserServiceMock)
				mockSvc.On("UpdatePassword", mock.Anything).
					Return(&dto.UpdatePasswordResponse{
						Message: "Password updated successfully",
					}, nil)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			expectedBody: dto.UpdatePasswordResponse{
				Message: "Password updated successfully",
			},
		},
		{
			name:           "empty body (EOF)",
			requestBody:    ``,
			setupMock:      func() *user_service.UserServiceMock { return new(user_service.UserServiceMock) },
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Request body cannot be empty",
			},
		},
		{
			name: "missing fields (validation error)",
			requestBody: map[string]string{
				"old_password": "",
			},
			setupMock:      func() *user_service.UserServiceMock { return new(user_service.UserServiceMock) },
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Key: 'UpdatePasswordRequest.ID' Error:Field validation for 'ID' failed on the 'required' tag\nKey: 'UpdatePasswordRequest.OldPassword' Error:Field validation for 'OldPassword' failed on the 'required' tag\nKey: 'UpdatePasswordRequest.NewPassword' Error:Field validation for 'NewPassword' failed on the 'required' tag",
			},
		},

		{
			name: "userID missing ",
			requestBody: dto.UpdatePasswordRequest{
				OldPassword: "superoldpassword",
				NewPassword: "ReallyStrongPass123!",
			},
			setupMock:      func() *user_service.UserServiceMock { return new(user_service.UserServiceMock) },
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Key: 'UpdatePasswordRequest.ID' Error:Field validation for 'ID' failed on the 'required' tag",
			},
		},

		{
			name: "record not found this",
			requestBody: dto.UpdatePasswordRequest{
				ID:          1,
				OldPassword: "realloldpassword!",
				NewPassword: "ReallyStrongPass123s!",
			},
			setupMock: func() *user_service.UserServiceMock {

				mockSvc := new(user_service.UserServiceMock)
				mockSvc.On("UpdatePassword", mock.Anything).
					Return(nil, gorm.ErrRecordNotFound)
				return mockSvc
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: common.ErrorResponse{
				Message: "user not found",
			},
		},

		{
			name: "invalid old password",
			requestBody: dto.UpdatePasswordRequest{
				ID:          1,
				OldPassword: "reallyoldpassword!",
				NewPassword: "invalidoldpassword",
			},
			setupMock: func() *user_service.UserServiceMock {
				mockSvc := new(user_service.UserServiceMock)
				mockSvc.On("UpdatePassword", mock.Anything).Return(nil, errors.New("invalid old password"))
				return mockSvc
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "invalid old password",
			},
		},

		{
			name: "generic error",
			requestBody: dto.UpdatePasswordRequest{
				ID:          1,
				OldPassword: "reallyoldpassword!",
				NewPassword: "somethingfailed",
			},
			setupMock: func() *user_service.UserServiceMock {

				mockSvc := new(user_service.UserServiceMock)
				mockSvc.On("UpdatePassword", mock.Anything).Return(nil, errors.New("something failed"))

				return mockSvc
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "something failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()

			handler := NewUserHandler(mockService, nil)
			router := setupGin()

			// inject fake userID (important)
			router.Use(func(c *gin.Context) {
				c.Set("userID", 1)
				c.Next()
			})

			router.PATCH("/users/password", handler.UpdatePassword)

			var body []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPatch, "/users/password", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var actual any
			_ = json.Unmarshal(w.Body.Bytes(), &actual)

			expectedBytes, _ := json.Marshal(tt.expectedBody)
			var expected any
			_ = json.Unmarshal(expectedBytes, &expected)

			assert.Equal(t, expected, actual)

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_UpdatePassword_Context_Failures(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(c *gin.Context)
		requestBody    any
		setupMock      func() *user_service.UserServiceMock
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "missing userID in context -> 401",
			setupContext: func(c *gin.Context) {
			},
			requestBody: dto.UpdatePasswordRequest{
				ID:          1,
				OldPassword: "superoldpassword",
				NewPassword: "ReallyStrongPass123!",
			},
			setupMock: func() *user_service.UserServiceMock {
				return new(user_service.UserServiceMock)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   common.ErrorResponse{Message: "unauthorized"},
		},
		{
			name: "userID is wrong type -> panic",
			setupContext: func(c *gin.Context) {
				c.Set("userID", "not-an-int")
			},
			requestBody: dto.UpdatePasswordRequest{
				ID:          1,
				OldPassword: "superoldpassword",
				NewPassword: "ReallyStrongPass123!",
			},
			setupMock: func() *user_service.UserServiceMock {
				return new(user_service.UserServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   common.ErrorResponse{Message: "invalid userID in context"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()

			handler := NewUserHandler(mockService, nil)
			router := setupGin()

			router.Use(func(c *gin.Context) {
				tt.setupContext(c)
				c.Next()
			})
			router.PATCH("/users/password", handler.UpdatePassword)

			var body []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPatch, "/users/password", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var actual any
			_ = json.Unmarshal(w.Body.Bytes(), &actual)

			expectedBytes, _ := json.Marshal(tt.expectedBody)
			var expected any
			_ = json.Unmarshal(expectedBytes, &expected)

			assert.Equal(t, expected, actual)

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(c *gin.Context)
		setupMock      func() *user_service.UserServiceMock
		expectedStatus int
		expectedBody   any
	}{

		{
			name:         "success case",
			setupContext: func(c *gin.Context) { c.Set("userID", 1) },
			setupMock: func() *user_service.UserServiceMock {
				mockSvc := new(user_service.UserServiceMock)
				mockSvc.On("DeleteUser", mock.Anything).
					Return(&dto.DeleteUserResponse{
						Message: "User account deleted successfully",
					}, nil)
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			expectedBody: dto.DeleteUserResponse{
				Message: "User account deleted successfully",
			},
		},
		{
			name: "missing userID in context -> 401",
			setupContext: func(c *gin.Context) {
			},
			setupMock: func() *user_service.UserServiceMock {
				return new(user_service.UserServiceMock)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   common.ErrorResponse{Message: "unauthorized"},
		},

		{
			name:         "user not found in service -> 404",
			setupContext: func(c *gin.Context) { c.Set("userID", 1) },
			setupMock: func() *user_service.UserServiceMock {
				m := new(user_service.UserServiceMock)
				m.On("DeleteUser", mock.Anything).
					Return(nil, gorm.ErrRecordNotFound)
				return m
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   common.ErrorResponse{Message: "user not found"},
		},
		{
			name: "userID is wrong type -> panic",
			setupContext: func(c *gin.Context) {
				c.Set("userID", "not-an-int")
			},
			setupMock: func() *user_service.UserServiceMock {
				return new(user_service.UserServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   common.ErrorResponse{Message: "invalid userID in context"},
		},
		{
			name:         "service returns generic error -> 400",
			setupContext: func(c *gin.Context) { c.Set("userID", 1) },
			setupMock: func() *user_service.UserServiceMock {
				m := new(user_service.UserServiceMock)
				m.On("DeleteUser", mock.Anything).
					Return(nil, errors.New("db failure"))
				return m
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   common.ErrorResponse{Message: "db failure"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := tt.setupMock()
			handler := NewUserHandler(mockSvc, nil)

			router := setupGin()

			router.Use(func(c *gin.Context) {
				tt.setupContext(c)
				c.Next()
			})

			router.DELETE("/users/account", handler.DeleteUser)

			req := httptest.NewRequest(http.MethodDelete, "/users/account", bytes.NewBuffer([]byte("{}")))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var actual any
			_ = json.Unmarshal(w.Body.Bytes(), &actual)

			exp, _ := json.Marshal(tt.expectedBody)
			var expected any
			_ = json.Unmarshal(exp, &expected)

			assert.Equal(t, expected, actual)

			mockSvc.AssertExpectations(t)
		})
	}
}

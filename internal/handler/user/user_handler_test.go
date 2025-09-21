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

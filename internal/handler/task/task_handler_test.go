package task_handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"taskflow/internal/auth"
	"taskflow/internal/common"
	"taskflow/internal/dto"
	task_service "taskflow/internal/service/task"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func setupGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestTaskHandler_CreateTask(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		requestBody    any
		setupMock      func() *task_service.TaskServiceMock
		expectedStatus int
		expectedBody   any
	}{
		{
			name:   "success case",
			userID: 1,
			requestBody: dto.CreateTaskRequest{
				Task: "Buy Milk",
			},
			setupMock: func() *task_service.TaskServiceMock {
				mockService := new(task_service.TaskServiceMock)
				mockService.On("CreateTask", 1, mock.MatchedBy(func(req *dto.CreateTaskRequest) bool {
					return req.Task == "Buy Milk"
				})).Return(nil)
				return mockService
			},
			expectedStatus: http.StatusCreated,
			expectedBody: dto.CreateTaskRequest{
				Task: "Buy Milk",
			},
		},
		{
			name:   "failure case - no userID",
			userID: 0,
			requestBody: dto.CreateTaskRequest{
				Task: "Buy Milk",
			},
			setupMock: func() *task_service.TaskServiceMock {
				return new(task_service.TaskServiceMock)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: common.ErrorResponse{
				Message: "unauthorized",
			},
		},
		{
			name:        "failure case - invalid JSON",
			userID:      1,
			requestBody: `{"task": }`, // malformed JSON
			setupMock: func() *task_service.TaskServiceMock {
				return new(task_service.TaskServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "invalid character '}' looking for beginning of value",
			},
		},
		{
			name:   "failure case - service error",
			userID: 1,
			requestBody: dto.CreateTaskRequest{
				Task: "Buy Milk",
			},
			setupMock: func() *task_service.TaskServiceMock {
				mockService := new(task_service.TaskServiceMock)
				mockService.On("CreateTask", 1, mock.Anything).Return(errors.New("service error"))
				return mockService
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "service error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			mockAuth := new(auth.MockUserAuth)
			handler := NewTaskHandler(mockService, mockAuth)

			router := setupGin()
			router.POST("/tasks", func(c *gin.Context) {
				if tt.userID != 0 {
					c.Set("userID", tt.userID) // inject userID to simulate authentication
				}
				handler.CreateTask(c)
			})

			var body []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody any
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedBodyBytes, err := json.Marshal(tt.expectedBody)
			assert.NoError(t, err)

			var expectedResponseBody any
			err = json.Unmarshal(expectedBodyBytes, &expectedResponseBody)
			assert.NoError(t, err)

			// For validation error, just check if it's an error response
			if tt.name == "failure case - validation error" {
				errorResp := make(map[string]any)
				json.Unmarshal(w.Body.Bytes(), &errorResp)
				assert.Contains(t, errorResp["error"], "required")
			} else {
				assert.Equal(t, expectedResponseBody, responseBody)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestTaskHandler_GetTask(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		taskID         string
		setupMock      func() *task_service.TaskServiceMock
		expectedStatus int
		expectedBody   any
	}{
		{
			name:   "success case",
			userID: 1,
			taskID: "1",
			setupMock: func() *task_service.TaskServiceMock {
				mockService := new(task_service.TaskServiceMock)
				mockService.On("GetTask", 1, 1).Return(dto.GetTaskResponse{
					ID:     1,
					Task:   "Buy Milk",
					Status: "pending",
				}, nil)
				return mockService
			},
			expectedStatus: http.StatusOK,
			expectedBody: dto.GetTaskResponse{
				ID:     1,
				Task:   "Buy Milk",
				Status: "pending",
			},
		},
		{
			name:   "failure case - invalid ID",
			userID: 1,
			taskID: "invalid",
			setupMock: func() *task_service.TaskServiceMock {
				return new(task_service.TaskServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Invalid ID",
			},
		},
		{
			name:   "failure case - ID less than 1",
			userID: 1,
			taskID: "0",
			setupMock: func() *task_service.TaskServiceMock {
				return new(task_service.TaskServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Invalid ID",
			},
		},
		{
			name:   "failure case - task not found",
			userID: 1,
			taskID: "999",
			setupMock: func() *task_service.TaskServiceMock {
				mockService := new(task_service.TaskServiceMock)
				mockService.On("GetTask", 1, 999).Return(dto.GetTaskResponse{}, errors.New("not found"))
				return mockService
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: common.ErrorResponse{
				Message: "Task not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			mockAuth := new(auth.MockUserAuth)
			handler := NewTaskHandler(mockService, mockAuth)

			router := setupGin()
			// Wrap handler to inject userID simulating authenticated request
			router.GET("/tasks/:id", func(c *gin.Context) {
				c.Set("userID", tt.userID)
				handler.GetTask(c)
			})

			req := httptest.NewRequest(http.MethodGet, "/tasks/"+tt.taskID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody any
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedBodyBytes, err := json.Marshal(tt.expectedBody)
			assert.NoError(t, err)

			var expectedResponseBody any
			err = json.Unmarshal(expectedBodyBytes, &expectedResponseBody)
			assert.NoError(t, err)

			assert.Equal(t, expectedResponseBody, responseBody)

			mockService.AssertExpectations(t)
		})
	}
}

func TestTaskHandler_ListTasks(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func() *task_service.TaskServiceMock
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "success case - with tasks",
			setupMock: func() *task_service.TaskServiceMock {
				mockService := new(task_service.TaskServiceMock)
				mockService.On("ListTasks").Return(dto.ListTasksResponse{
					Tasks: []dto.GetTaskResponse{
						{ID: 1, Task: "Buy Milk", Status: "pending"},
						{ID: 2, Task: "Buy Eggs", Status: "completed"},
					},
				}, nil)
				return mockService
			},
			expectedStatus: http.StatusOK,
			expectedBody: dto.ListTasksResponse{
				Tasks: []dto.GetTaskResponse{
					{ID: 1, Task: "Buy Milk", Status: "pending"},
					{ID: 2, Task: "Buy Eggs", Status: "completed"},
				},
			},
		},
		{
			name: "success case - empty list",
			setupMock: func() *task_service.TaskServiceMock {
				mockService := new(task_service.TaskServiceMock)
				mockService.On("ListTasks").Return(dto.ListTasksResponse{
					Tasks: []dto.GetTaskResponse{},
				}, nil)
				return mockService
			},
			expectedStatus: http.StatusOK,
			expectedBody: dto.ListTasksResponse{
				Tasks: []dto.GetTaskResponse{},
			},
		},
		{
			name: "failure case - service error",
			setupMock: func() *task_service.TaskServiceMock {
				mockService := new(task_service.TaskServiceMock)
				mockService.On("ListTasks").Return(dto.ListTasksResponse{}, errors.New("database error"))
				return mockService
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: common.ErrorResponse{
				Message: "database error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			mockAuth := new(auth.MockUserAuth)
			handler := NewTaskHandler(mockService, mockAuth)

			router := setupGin()
			router.GET("/tasks", handler.ListTasks)

			req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody any
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedBodyBytes, err := json.Marshal(tt.expectedBody)
			assert.NoError(t, err)

			var expectedResponseBody any
			err = json.Unmarshal(expectedBodyBytes, &expectedResponseBody)
			assert.NoError(t, err)

			assert.Equal(t, expectedResponseBody, responseBody)

			mockService.AssertExpectations(t)
		})
	}
}

func TestTaskHandler_UpdateStatus(t *testing.T) {
	tests := []struct {
		name           string
		taskID         string
		requestBody    any
		setupMock      func() *task_service.TaskServiceMock
		expectedStatus int
		expectedBody   any
	}{
		{
			name:   "success case",
			taskID: "1",
			requestBody: dto.UpdateStatusRequest{
				Status: "completed",
			},
			setupMock: func() *task_service.TaskServiceMock {
				mockService := new(task_service.TaskServiceMock)
				mockService.On("UpdateStatus", 1, "completed").Return(nil)
				return mockService
			},
			expectedStatus: http.StatusOK,
			expectedBody: dto.UpdateStatusResponse{
				Message: "status updated",
			},
		},
		{
			name:   "failure case - invalid ID",
			taskID: "invalid",
			requestBody: dto.UpdateStatusRequest{
				Status: "completed",
			},
			setupMock: func() *task_service.TaskServiceMock {
				return new(task_service.TaskServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "invalid task ID",
			},
		},
		{
			name:   "failure case - ID less than 1",
			taskID: "0",
			requestBody: dto.UpdateStatusRequest{
				Status: "completed",
			},
			setupMock: func() *task_service.TaskServiceMock {
				return new(task_service.TaskServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "invalid task ID",
			},
		},
		{
			name:        "failure case - invalid JSON",
			taskID:      "1",
			requestBody: `{"status": }`, // malformed JSON
			setupMock: func() *task_service.TaskServiceMock {
				return new(task_service.TaskServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "invalid character '}' looking for beginning of value",
			},
		},
		{
			name:   "failure case - invalid status",
			taskID: "1",
			requestBody: dto.UpdateStatusRequest{
				Status: "invalid-status",
			},
			setupMock: func() *task_service.TaskServiceMock {
				return new(task_service.TaskServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "failure case - task not found",
			taskID: "999",
			requestBody: dto.UpdateStatusRequest{
				Status: "completed",
			},
			setupMock: func() *task_service.TaskServiceMock {
				mockService := new(task_service.TaskServiceMock)
				mockService.On("UpdateStatus", 999, "completed").Return(gorm.ErrRecordNotFound)
				return mockService
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: common.ErrorResponse{
				Message: "Task not found",
			},
		},
		{
			name:   "failure case - service error",
			taskID: "1",
			requestBody: dto.UpdateStatusRequest{
				Status: "completed",
			},
			setupMock: func() *task_service.TaskServiceMock {
				mockService := new(task_service.TaskServiceMock)
				mockService.On("UpdateStatus", 1, "completed").Return(errors.New("service error"))
				return mockService
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "service error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			mockAuth := new(auth.MockUserAuth)
			handler := NewTaskHandler(mockService, mockAuth)

			router := setupGin()
			router.PATCH("/tasks/:id/status", handler.UpdateStatus)

			var body []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPatch, "/tasks/"+tt.taskID+"/status", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var responseBody any
				err = json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				// Special handling for validation errors
				if tt.name == "failure case - invalid status" {
					errorResp := make(map[string]any)
					json.Unmarshal(w.Body.Bytes(), &errorResp)
					assert.Contains(t, errorResp["error"], "oneof")
				} else {
					expectedBodyBytes, err := json.Marshal(tt.expectedBody)
					assert.NoError(t, err)

					var expectedResponseBody any
					err = json.Unmarshal(expectedBodyBytes, &expectedResponseBody)
					assert.NoError(t, err)

					assert.Equal(t, expectedResponseBody, responseBody)
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestTaskHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		taskID         string
		setupMock      func() *task_service.TaskServiceMock
		expectedStatus int
		expectedBody   any
	}{
		{
			name:   "success case",
			taskID: "1",
			setupMock: func() *task_service.TaskServiceMock {
				mockService := new(task_service.TaskServiceMock)
				mockService.On("Delete", 1).Return(nil)
				return mockService
			},
			expectedStatus: http.StatusOK,
			expectedBody: dto.DeleteTaskResponse{
				Message: "Task deleted successfully",
			},
		},
		{
			name:   "failure case - invalid ID",
			taskID: "invalid",
			setupMock: func() *task_service.TaskServiceMock {
				return new(task_service.TaskServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Invalid ID",
			},
		},
		{
			name:   "failure case - ID less than 1",
			taskID: "0",
			setupMock: func() *task_service.TaskServiceMock {
				return new(task_service.TaskServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Invalid ID",
			},
		},
		{
			name:   "failure case - task not found",
			taskID: "999",
			setupMock: func() *task_service.TaskServiceMock {
				mockService := new(task_service.TaskServiceMock)
				mockService.On("Delete", 999).Return(gorm.ErrRecordNotFound)
				return mockService
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: common.ErrorResponse{
				Message: "Task not found",
			},
		},
		{
			name:   "failure case - service error",
			taskID: "1",
			setupMock: func() *task_service.TaskServiceMock {
				mockService := new(task_service.TaskServiceMock)
				mockService.On("Delete", 1).Return(errors.New("database error"))
				return mockService
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: common.ErrorResponse{
				Message: "Couldn't delete task",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			mockAuth := new(auth.MockUserAuth)
			handler := NewTaskHandler(mockService, mockAuth)

			router := setupGin()
			router.DELETE("/tasks/:id", handler.Delete)

			req := httptest.NewRequest(http.MethodDelete, "/tasks/"+tt.taskID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody any
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			expectedBodyBytes, err := json.Marshal(tt.expectedBody)
			assert.NoError(t, err)

			var expectedResponseBody any
			err = json.Unmarshal(expectedBodyBytes, &expectedResponseBody)
			assert.NoError(t, err)

			assert.Equal(t, expectedResponseBody, responseBody)

			mockService.AssertExpectations(t)
		})
	}
}

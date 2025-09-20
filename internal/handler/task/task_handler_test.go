package task_handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"taskflow/internal/common"
	"taskflow/internal/dto"
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
		requestBody    any
		setupMock      func() *TaskServiceMock
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "success case",
			requestBody: dto.CreateTaskRequest{
				Task: "Buy Milk",
			},
			setupMock: func() *TaskServiceMock {
				mockService := new(TaskServiceMock)
				mockService.On("CreateTask", mock.MatchedBy(func(req *dto.CreateTaskRequest) bool {
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
			name:        "failure case - invalid JSON",
			requestBody: `{"task": }`, // malformed JSON
			setupMock: func() *TaskServiceMock {
				return new(TaskServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "invalid character '}' looking for beginning of value",
			},
		},
		{
			name: "failure case - validation error",
			requestBody: dto.CreateTaskRequest{
				Task: "", // empty task
			},
			setupMock: func() *TaskServiceMock {
				return new(TaskServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Key: 'CreateTaskRequest.Task' Error:Tag: 'required' ActualTag: 'required' Namespace: 'CreateTaskRequest.Task' StructNamespace: 'CreateTaskRequest.Task' StructField: 'Task' ActualField: 'Task' Value: '' Param: ''",
			},
		},
		{
			name: "failure case - service error",
			requestBody: dto.CreateTaskRequest{
				Task: "Buy Milk",
			},
			setupMock: func() *TaskServiceMock {
				mockService := new(TaskServiceMock)
				mockService.On("CreateTask", mock.Anything).Return(errors.New("service error"))
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
			handler := NewTaskHandler(mockService)

			router := setupGin()
			router.POST("/tasks", handler.CreateTask)

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
		taskID         string
		setupMock      func() *TaskServiceMock
		expectedStatus int
		expectedBody   any
	}{
		{
			name:   "success case",
			taskID: "1",
			setupMock: func() *TaskServiceMock {
				mockService := new(TaskServiceMock)
				mockService.On("GetTask", 1).Return(dto.GetTaskResponse{
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
			taskID: "invalid",
			setupMock: func() *TaskServiceMock {
				return new(TaskServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Invalid ID",
			},
		},
		{
			name:   "failure case - ID less than 1",
			taskID: "0",
			setupMock: func() *TaskServiceMock {
				return new(TaskServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Invalid ID",
			},
		},
		{
			name:   "failure case - task not found",
			taskID: "999",
			setupMock: func() *TaskServiceMock {
				mockService := new(TaskServiceMock)
				mockService.On("GetTask", 999).Return(dto.GetTaskResponse{}, errors.New("not found"))
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
			handler := NewTaskHandler(mockService)

			router := setupGin()
			router.GET("/tasks/:id", handler.GetTask)

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
		setupMock      func() *TaskServiceMock
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "success case - with tasks",
			setupMock: func() *TaskServiceMock {
				mockService := new(TaskServiceMock)
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
			setupMock: func() *TaskServiceMock {
				mockService := new(TaskServiceMock)
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
			setupMock: func() *TaskServiceMock {
				mockService := new(TaskServiceMock)
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
			handler := NewTaskHandler(mockService)

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
		setupMock      func() *TaskServiceMock
		expectedStatus int
		expectedBody   any
	}{
		{
			name:   "success case",
			taskID: "1",
			requestBody: dto.UpdateStatusRequest{
				Status: "completed",
			},
			setupMock: func() *TaskServiceMock {
				mockService := new(TaskServiceMock)
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
			setupMock: func() *TaskServiceMock {
				return new(TaskServiceMock)
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
			setupMock: func() *TaskServiceMock {
				return new(TaskServiceMock)
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
			setupMock: func() *TaskServiceMock {
				return new(TaskServiceMock)
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
			setupMock: func() *TaskServiceMock {
				return new(TaskServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "failure case - task not found",
			taskID: "999",
			requestBody: dto.UpdateStatusRequest{
				Status: "completed",
			},
			setupMock: func() *TaskServiceMock {
				mockService := new(TaskServiceMock)
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
			setupMock: func() *TaskServiceMock {
				mockService := new(TaskServiceMock)
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
			handler := NewTaskHandler(mockService)

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
		setupMock      func() *TaskServiceMock
		expectedStatus int
		expectedBody   any
	}{
		{
			name:   "success case",
			taskID: "1",
			setupMock: func() *TaskServiceMock {
				mockService := new(TaskServiceMock)
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
			setupMock: func() *TaskServiceMock {
				return new(TaskServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Invalid ID",
			},
		},
		{
			name:   "failure case - ID less than 1",
			taskID: "0",
			setupMock: func() *TaskServiceMock {
				return new(TaskServiceMock)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: common.ErrorResponse{
				Message: "Invalid ID",
			},
		},
		{
			name:   "failure case - task not found",
			taskID: "999",
			setupMock: func() *TaskServiceMock {
				mockService := new(TaskServiceMock)
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
			setupMock: func() *TaskServiceMock {
				mockService := new(TaskServiceMock)
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
			handler := NewTaskHandler(mockService)

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

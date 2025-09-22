package task_service

import (
	"errors"
	"taskflow/internal/domain/task"
	"taskflow/internal/dto"
	"taskflow/internal/repository/gorm/gorm_task"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTaskService_CreateTask(t *testing.T) {
	tests := []struct {
		name        string
		userID      int
		taskRequest *dto.CreateTaskRequest
		setupMock   func() *gorm_task.TaskRepoMock
		wantErr     bool
		errMessage  string
	}{
		{
			name:   "success case - create task Buy Milk",
			userID: 1,
			taskRequest: &dto.CreateTaskRequest{
				Task: "Buy Milk",
			},
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("Create", mock.MatchedBy(func(tk *task.Task) bool {
					return tk.UserID == 1 && tk.Task == "Buy Milk" && tk.Status == "pending"
				})).Return(nil)
				return mockRepo
			},
			wantErr:    false,
			errMessage: "",
		},
		{
			name:   "failure case - empty task",
			userID: 1,
			taskRequest: &dto.CreateTaskRequest{
				Task: "",
			},
			setupMock: func() *gorm_task.TaskRepoMock {
				return new(gorm_task.TaskRepoMock)
			},
			wantErr:    true,
			errMessage: "task name cannot be empty",
		},
		{
			name: "failure case - invalid user",
			taskRequest: &dto.CreateTaskRequest{
				Task: "Buy Milk", // <-- non-empty task
			},
			setupMock: func() *gorm_task.TaskRepoMock {
				return new(gorm_task.TaskRepoMock)
			},
			userID:     0, // <-- invalid user triggers the error
			wantErr:    true,
			errMessage: "invalid user",
		},
		{
			name:   "failure case - database error",
			userID: 2,
			taskRequest: &dto.CreateTaskRequest{
				Task: "Buy Eggs",
			},
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("Create", mock.MatchedBy(func(tk *task.Task) bool {
					return tk.UserID == 2 && tk.Task == "Buy Eggs" && tk.Status == "pending"
				})).Return(errors.New("db error"))
				return mockRepo
			},
			wantErr:    true,
			errMessage: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			service := NewTaskService(mockRepo)

			err := service.CreateTask(tt.userID, tt.taskRequest)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.errMessage)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTaskService_GetTask(t *testing.T) {
	tests := []struct {
		name      string // description of this test case
		id        int
		userID    int
		setupMock func() *gorm_task.TaskRepoMock
		want      dto.GetTaskResponse
		wantErr   bool
	}{
		{
			name:   "success case",
			id:     1,
			userID: 1,
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("GetByID", 1, 1).Return(&task.Task{
					ID:     1,
					Task:   "Buy Milk",
					Status: "pending",
					UserID: 1,
				}, nil)
				return mockRepo
			},
			want: dto.GetTaskResponse{
				ID:     1,
				Task:   "Buy Milk",
				Status: "pending",
			},
			wantErr: false,
		},
		{
			name:   "failure case - task not found",
			id:     2,
			userID: 1,
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("GetByID", 1, 2).Return((*task.Task)(nil), errors.New("not found"))
				return mockRepo
			},
			want:    dto.GetTaskResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			s := NewTaskService(mockRepo)
			got, gotErr := s.GetTask(tt.userID, tt.id)

			if tt.wantErr {
				assert.Error(t, gotErr)
				assert.Equal(t, dto.GetTaskResponse{}, got)
			} else {
				assert.NoError(t, gotErr)
				assert.Equal(t, tt.want, got)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTaskService_ListTasks(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		setupMock func() *gorm_task.TaskRepoMock
		want      dto.ListTasksResponse
		wantErr   bool
	}{
		{
			name:   "success - single task",
			userID: 1,
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("List", 1).Return([]task.Task{
					{ID: 1, Task: "Buy milk", Status: "pending", CreatedAt: time.Now()},
				}, nil)
				return mockRepo
			},
			want: dto.ListTasksResponse{
				Tasks: []dto.GetTaskResponse{
					{ID: 1, Task: "Buy milk", Status: "pending"},
				},
			},
			wantErr: false,
		},
		{
			name:   "failure - invalid userID",
			userID: 0,
			setupMock: func() *gorm_task.TaskRepoMock {
				return new(gorm_task.TaskRepoMock) // repo should not be called
			},
			want:    dto.ListTasksResponse{},
			wantErr: true,
		},
		{
			name:   "failure - db error",
			userID: 1,
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				// return empty slice instead of nil
				mockRepo.On("List", 1).Return([]task.Task{}, errors.New("db error"))
				return mockRepo
			},
			want:    dto.ListTasksResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			s := NewTaskService(mockRepo)
			got, gotErr := s.ListTasks(tt.userID)

			if tt.wantErr {
				assert.Error(t, gotErr)
				assert.Equal(t, dto.ListTasksResponse{}, got)
			} else {
				assert.NoError(t, gotErr)
				assert.Equal(t, tt.want, got)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTaskService_UpdateStatus(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func() *gorm_task.TaskRepoMock
		userID    int
		id        int
		status    string
		wantErr   bool
	}{
		{
			name:   "success - pending",
			userID: 123,
			id:     1,
			status: "pending",
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("UpdateStatus",
					123,
					1,
					"pending",
				).Return(nil)
				return mockRepo
			},
			wantErr: false,
		},
		{
			name:   "success - completed",
			userID: 123,
			id:     2,
			status: "completed",
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("UpdateStatus", 123, 2, "completed").Return(nil)
				return mockRepo
			},
			wantErr: false,
		},
		{
			name:   "failure - invalid status",
			userID: 123,
			id:     3,
			status: "invalid",
			setupMock: func() *gorm_task.TaskRepoMock {
				return new(gorm_task.TaskRepoMock)
			},
			wantErr: true,
		},
		{
			name:   "failure - repo error",
			userID: 123,
			id:     4,
			status: "pending",
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("UpdateStatus", 123, 4, "pending").Return(errors.New("db error"))
				return mockRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			s := NewTaskService(mockRepo)

			gotErr := s.UpdateStatus(tt.userID, tt.id, tt.status)

			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTaskService_Delete(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func() *gorm_task.TaskRepoMock
		userID    int
		id        int
		wantErr   bool
	}{
		{
			name:   "success",
			userID: 1,
			id:     1,
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("Delete", 1, 1).Return(nil)
				return mockRepo
			},
			wantErr: false,
		},
		{
			name:   "failure - repo error",
			userID: 1,
			id:     2,
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("Delete", 1, 2).Return(errors.New("db error"))
				return mockRepo
			},
			wantErr: true,
		},
		{
			name:   "failure - delete non-existing task",
			userID: 1,
			id:     3,
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("Delete", 1, 3).Return(errors.New("not found"))
				return mockRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			s := NewTaskService(mockRepo)

			err := s.Delete(tt.userID, tt.id)

			if tt.wantErr {
				assert.Error(t, err, "expected an error but got nil")
			} else {
				assert.NoError(t, err, "expected no error but got one")
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

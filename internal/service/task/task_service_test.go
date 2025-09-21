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
		taskRequest *dto.CreateTaskRequest
		setupMock   func() *gorm_task.TaskRepoMock
		wantErr     bool
	}{
		{
			name: "success case - create task Buy Milk",
			taskRequest: &dto.CreateTaskRequest{
				Task: "Buy Milk",
			},
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("Create", mock.MatchedBy(func(tk *task.Task) bool {
					return tk.Task == "Buy Milk"
				})).Return(nil)
				return mockRepo
			},
			wantErr: false,
		},
		{
			name: "failure case - Empty Task",
			taskRequest: &dto.CreateTaskRequest{
				Task: "",
			},
			setupMock: func() *gorm_task.TaskRepoMock {
				return new(gorm_task.TaskRepoMock)
			},
			wantErr: true,
		},
		{
			name: "failure case - database error",
			taskRequest: &dto.CreateTaskRequest{
				Task: "Buy Eggs",
			},
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("Create", mock.MatchedBy(func(tk *task.Task) bool {
					return tk.Task == "Buy Eggs"
				})).Return(errors.New("db error"))
				return mockRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			service := NewTaskService(mockRepo)

			err := service.CreateTask(tt.taskRequest)

			if tt.wantErr {
				assert.Error(t, err)
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
		setupMock func() *gorm_task.TaskRepoMock
		want      dto.GetTaskResponse
		wantErr   bool
	}{
		{
			name: "success case",
			id:   1,
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("GetByID", 1).Return(&task.Task{
					ID:     1,
					Task:   "Buy Milk",
					Status: "pending",
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
			name: "failure case - task not found",
			id:   2,
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("GetByID", 2).Return((*task.Task)(nil), errors.New("not found"))
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
			got, gotErr := s.GetTask(tt.id)

			if tt.wantErr {
				assert.Error(t, gotErr)
				assert.Equal(t, dto.GetTaskResponse{}, got)
			} else {
				assert.NotZero(t, got)
				assert.NoError(t, gotErr)
				assert.Equal(t, tt.want, got)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTaskService_ListTasks(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		setupMock func() *gorm_task.TaskRepoMock
		want      dto.ListTasksResponse
		wantErr   bool
	}{
		{
			name: "success",
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("List").Return([]task.Task{
					{
						ID: 1, Task: "Buy milk", Status: "pending", CreatedAt: time.Now(),
					},
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
			name: "success - multiple tasks",
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("List").Return([]task.Task{
					{ID: 1, Task: "Buy Milk", Status: "pending", CreatedAt: time.Now()},
					{ID: 2, Task: "Buy Eggs", Status: "completed", CreatedAt: time.Now()},
				}, nil)
				return mockRepo
			},
			want: dto.ListTasksResponse{
				Tasks: []dto.GetTaskResponse{
					{ID: 1, Task: "Buy Milk", Status: "pending"},
					{ID: 2, Task: "Buy Eggs", Status: "completed"},
				},
			},
			wantErr: false,
		},
		{
			name: "failure - db error",
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("List").Return(([]task.Task)(nil), errors.New("db error"))
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
			got, gotErr := s.ListTasks()

			if tt.wantErr {
				assert.Error(t, gotErr)
				assert.Equal(t, dto.ListTasksResponse{}, got)
			} else {
				assert.NotZero(t, got)
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
		id        int
		status    string
		wantErr   bool
	}{
		{
			name:   "success - pending",
			id:     1,
			status: "pending",
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("UpdateStatus",
					mock.MatchedBy(func(id int) bool { return id > 0 }),
					mock.MatchedBy(func(status string) bool { return status == "pending" || status == "completed" }),
				).Return(nil)
				return mockRepo
			},
			wantErr: false,
		},
		{
			name:   "success - completed",
			id:     2,
			status: "completed",
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("UpdateStatus",
					mock.Anything,
					mock.Anything,
				).Return(nil)
				return mockRepo
			},
			wantErr: false,
		},
		{
			name:   "failure - invalid status",
			id:     3,
			status: "invalid",
			setupMock: func() *gorm_task.TaskRepoMock {
				return new(gorm_task.TaskRepoMock)
			},
			wantErr: true,
		},
		{
			name:   "failure - repo error",
			id:     4,
			status: "pending",
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("UpdateStatus", mock.Anything, mock.Anything).Return(errors.New("db error"))
				return mockRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			s := NewTaskService(mockRepo)
			gotErr := s.UpdateStatus(tt.id, tt.status)

			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}

			mockRepo.AssertExpectations(t)

			if tt.status == "pending" || tt.status == "completed" {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestTaskService_Delete(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func() *gorm_task.TaskRepoMock
		id        int
		wantErr   bool
	}{
		{
			name: "success",
			id:   1,
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("Delete", 1).Return(nil)
				return mockRepo
			},
			wantErr: false,
		},

		{
			name: "failure - repo error",
			id:   2,
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("Delete", 2).Return(errors.New("db error"))
				return mockRepo
			},
			wantErr: true,
		},

		{
			name: "failure - delete non-existing task",
			id:   3,
			setupMock: func() *gorm_task.TaskRepoMock {
				mockRepo := new(gorm_task.TaskRepoMock)
				mockRepo.On("Delete", 3).Return(errors.New("not found"))
				return mockRepo
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			s := NewTaskService(mockRepo)
			err := s.Delete(tt.id)

			assert.Equal(t, tt.wantErr, err != nil, "error mismatch")

			assert.True(t, mockRepo.AssertExpectations(t))
		})
	}
}

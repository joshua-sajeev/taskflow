package task_service

import (
	"taskflow/internal/dto"

	"github.com/stretchr/testify/mock"
)

type TaskServiceMock struct {
	mock.Mock
}

var _ TaskServiceInterface = (*TaskServiceMock)(nil)

func (m *TaskServiceMock) CreateTask(userID int, taskRequest *dto.CreateTaskRequest) error {

	args := m.Called(userID, taskRequest)
	return args.Error(0)
}
func (m *TaskServiceMock) GetTask(userID int, id int) (dto.GetTaskResponse, error) {
	args := m.Called(userID, id)
	return args.Get(0).(dto.GetTaskResponse), args.Error(1)
}
func (m *TaskServiceMock) ListTasks(userID int) (dto.ListTasksResponse, error) {
	args := m.Called(userID)
	return args.Get(0).(dto.ListTasksResponse), args.Error(1)
}
func (m *TaskServiceMock) UpdateStatus(userID int, id int, status string) error {
	args := m.Called(userID, id, status)
	return args.Error(0)
}
func (m *TaskServiceMock) Delete(userID int, id int) error {
	args := m.Called(userID, id)
	return args.Error(0)
}

package task_handler

import (
	"taskflow/internal/dto"
	task_service "taskflow/internal/service/task"

	"github.com/stretchr/testify/mock"
)

type TaskServiceMock struct {
	mock.Mock
}

var _ task_service.TaskServiceInterface = (*TaskServiceMock)(nil)

func (m *TaskServiceMock) CreateTask(taskRequest *dto.CreateTaskRequest) error {

	args := m.Called(taskRequest)
	return args.Error(0)
}
func (m *TaskServiceMock) GetTask(id int) (dto.GetTaskResponse, error) {
	args := m.Called(id)
	return args.Get(0).(dto.GetTaskResponse), args.Error(1)
}
func (m *TaskServiceMock) ListTasks() (dto.ListTasksResponse, error) {
	args := m.Called()
	return args.Get(0).(dto.ListTasksResponse), args.Error(1)
}
func (m *TaskServiceMock) UpdateStatus(id int, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}
func (m *TaskServiceMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

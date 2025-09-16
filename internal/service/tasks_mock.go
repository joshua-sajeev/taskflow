package service

import (
	"taskflow/internal/domain/task"

	"github.com/stretchr/testify/mock"
	gg "taskflow/internal/repository/gorm/gorm_task"
)

type TaskRepoMock struct {
	mock.Mock
}

var _ gg.TaskRepositoryInterface = (*TaskRepoMock)(nil)

func (m *TaskRepoMock) Create(task *task.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

// If panicking
//
//	if args.Get(0) == nil {
//		return nil, args.Error(1)
//	}
func (m *TaskRepoMock) GetByID(id int) (*task.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*task.Task), args.Error(1)
}

func (m *TaskRepoMock) List() ([]task.Task, error) {
	args := m.Called()
	return args.Get(0).([]task.Task), args.Error(1)
}

func (m *TaskRepoMock) Update(task *task.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *TaskRepoMock) Delete(id int) error {

	args := m.Called(id)
	return args.Error(0)
}

func (m *TaskRepoMock) UpdateStatus(id int, status string) error {

	args := m.Called(id, status)
	return args.Error(0)
}

package service

import (
	"errors"
	"taskflow/internal/domain/task"
)

type TaskService struct {
	repo task.TaskRepository
}

func NewTaskService(repo task.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(t *task.Task) error {
	if t.Task == "" {
		return errors.New("task name cannot be empty")
	}
	t.Status = "pending"
	return s.repo.Create(t)
}

func (s *TaskService) GetTask(id int) (*task.Task, error) {
	return s.repo.GetByID(id)
}

func (s *TaskService) ListTasks() ([]task.Task, error) {
	return s.repo.List()
}

func (s *TaskService) UpdateStatus(id int, status string) error {
	if status != "pending" && status != "completed" {
		return errors.New("invalid status")
	}
	return s.repo.UpdateStatus(id, status)
}

func (s *TaskService) Delete(id int) error {
	return s.repo.Delete(id)
}

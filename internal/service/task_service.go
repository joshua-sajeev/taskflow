package service

import (
	"errors"
	"taskflow/internal/domain/task"
	"taskflow/internal/dto"
)

type TaskService struct {
	repo task.TaskRepositoryInterface
}

func NewTaskService(repo task.TaskRepositoryInterface) *TaskService {
	return &TaskService{repo: repo}
}

var _ TaskServiceInterface = (*TaskService)(nil)

func (s *TaskService) CreateTask(taskRequest *dto.CreateTaskRequest) error {
	if taskRequest.Task == "" {
		return errors.New("task name cannot be empty")
	}

	task := task.Task{
		Task:   taskRequest.Task,
		Status: "pending",
	}

	return s.repo.Create(&task)
}
func (s *TaskService) GetTask(id int) (dto.GetTaskResponse, error) {
	t, err := s.repo.GetByID(id)
	if err != nil {
		return dto.GetTaskResponse{}, err
	}
	return dto.GetTaskResponse{
		ID:     t.ID,
		Task:   t.Task,
		Status: t.Status,
	}, nil
}

func (s *TaskService) ListTasks() (dto.ListTasksResponse, error) {
	tasks, err := s.repo.List() // This returns []task.Task
	if err != nil {
		return dto.ListTasksResponse{}, err
	}

	// Convert []task.Task to []dto.GetTaskResponse
	var taskResponses []dto.GetTaskResponse
	for _, task := range tasks {
		taskResponses = append(taskResponses, dto.GetTaskResponse{
			ID:     task.ID,
			Task:   task.Task,
			Status: task.Status,
		})
	}

	return dto.ListTasksResponse{
		Tasks: taskResponses,
	}, nil
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

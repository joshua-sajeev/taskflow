package task_service

import (
	"errors"
	"taskflow/internal/domain/task"
	"taskflow/internal/dto"
	"taskflow/internal/repository/gorm/gorm_task"
)

type TaskService struct {
	repo gorm_task.TaskRepositoryInterface
}

func NewTaskService(repo gorm_task.TaskRepositoryInterface) *TaskService {
	return &TaskService{repo: repo}
}

var _ TaskServiceInterface = (*TaskService)(nil)

func (s *TaskService) CreateTask(userID int, taskRequest *dto.CreateTaskRequest) error {
	if userID == 0 {
		return errors.New("invalid user")
	}

	if taskRequest.Task == "" {
		return errors.New("task name cannot be empty")
	}

	task := task.Task{
		UserID: userID,
		Task:   taskRequest.Task,
		Status: "pending",
	}

	return s.repo.Create(&task)
}

func (s *TaskService) GetTask(userID int, id int) (dto.GetTaskResponse, error) {
	if userID == 0 {
		return dto.GetTaskResponse{}, errors.New("invalid user")
	}

	t, err := s.repo.GetByID(userID, id)
	if err != nil {
		return dto.GetTaskResponse{}, err
	}
	return dto.GetTaskResponse{
		ID:     t.ID,
		Task:   t.Task,
		Status: t.Status,
	}, nil
}

func (s *TaskService) ListTasks(userID int) (dto.ListTasksResponse, error) {
	if userID == 0 {
		return dto.ListTasksResponse{}, errors.New("invalid user")
	}

	tasks, err := s.repo.List(userID) // This returns []task.Task
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

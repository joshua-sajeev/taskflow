package task_service

import (
	"taskflow/internal/dto"
)

type TaskServiceInterface interface {
	CreateTask(userID int, taskRequest *dto.CreateTaskRequest) error
	GetTask(userID int, id int) (dto.GetTaskResponse, error)
	ListTasks(userID int) (dto.ListTasksResponse, error)
	UpdateStatus(userID int, id int, status string) error
	Delete(userID int, id int) error
}

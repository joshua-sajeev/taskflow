package task_service

import "taskflow/internal/dto"

type TaskServiceInterface interface {
	CreateTask(userID int, taskRequest *dto.CreateTaskRequest) error
	GetTask(userID int, id int) (dto.GetTaskResponse, error)
	ListTasks() (dto.ListTasksResponse, error)
	UpdateStatus(id int, status string) error
	Delete(id int) error
}

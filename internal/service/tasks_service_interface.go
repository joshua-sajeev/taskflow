package service

import "taskflow/internal/dto"

type TaskServiceInterface interface {
	CreateTask(taskRequest *dto.CreateTaskRequest) error
	GetTask(id int) (dto.GetTaskResponse, error)
	ListTasks() (dto.ListTasksResponse, error)
	UpdateStatus(id int, status string) error
	Delete(id int) error
}

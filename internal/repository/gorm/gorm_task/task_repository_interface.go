package gorm_task

import "taskflow/internal/domain/task"

type TaskRepositoryInterface interface {
	Create(task *task.Task) error
	GetByID(id int) (*task.Task, error)
	List() ([]task.Task, error)
	Update(task *task.Task) error
	Delete(id int) error
	UpdateStatus(id int, status string) error
}

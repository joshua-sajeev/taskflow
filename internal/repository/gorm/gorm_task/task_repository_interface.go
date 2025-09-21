package gorm_task

import (
	"taskflow/internal/domain/task"
)

type TaskRepositoryInterface interface {
	Create(task *task.Task) error
	GetByID(userID int, id int) (*task.Task, error)
	List(userID int) ([]task.Task, error)
	Update(task *task.Task) error
	Delete(id int) error
	UpdateStatus(id int, status string) error
}

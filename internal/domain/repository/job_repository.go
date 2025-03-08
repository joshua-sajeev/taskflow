package repositories

import "taskflow/internal/domain/entities"

// JobRepository defines the methods for interacting with jobs in the database.
type JobRepository interface {
	Create(job *entities.Job) error
	FindByID(id uint) (*entities.Job, error)
	Update(job *entities.Job) error
	Delete(id uint) error
}

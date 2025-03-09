package repositories

import (
	"taskflow/internal/domain/entities"

	"github.com/google/uuid"
)

// JobRepository defines the methods for interacting with jobs in the database.
type JobRepository interface {
	Create(job *entities.Job) error
	FindByID(id uuid.UUID) (*entities.Job, error)
	Update(job *entities.Job) error
	Delete(id uuid.UUID) error
}

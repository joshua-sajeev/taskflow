package repositories

import (
	"taskflow/internal/domain/entities"

	"gorm.io/gorm"
)

// GormJobRepository is the GORM implementation of JobRepository.
type GormJobRepository struct {
	db *gorm.DB
}

// NewGormJobRepository creates a new instance of GormJobRepository.
func NewGormJobRepository(db *gorm.DB) JobRepository {
	return &GormJobRepository{db: db}
}

// Create inserts a new job into the database.
func (r *GormJobRepository) Create(job *entities.Job) error {
	return r.db.Create(job).Error
}

// FindByID retrieves a job by ID.
func (r *GormJobRepository) FindByID(id uint) (*entities.Job, error) {
	var job entities.Job
	err := r.db.First(&job, id).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// Update modifies an existing job.
func (r *GormJobRepository) Update(job *entities.Job) error {
	return r.db.Save(job).Error
}

// Delete removes a job by ID.
func (r *GormJobRepository) Delete(id uint) error {
	return r.db.Delete(&entities.Job{}, id).Error
}

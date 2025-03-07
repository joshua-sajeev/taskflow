package repositories

import (
	"taskflow/internal/domain/entities"

	"gorm.io/gorm"
)

type JobRepo struct {
	db *gorm.DB
}

type JobRepository interface {
	Create(job *entities.Job) error
	FindByID(id uint) (*entities.Job, error)
	Update(job *entities.Job) error
	Delete(id uint) error
}

func NewJobRepo(db *gorm.DB) *JobRepo {
	return &JobRepo{db: db}
}

func (r *JobRepo) Create(job *entities.Job) error {
	return r.db.Create(job).Error
}

func (r *JobRepo) FindByID(id uint) (*entities.Job, error) {
	var job entities.Job
	err := r.db.First(&job, id).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *JobRepo) Update(job *entities.Job) error {
	return r.db.Save(job).Error
}

func (r *JobRepo) Delete(id uint) error {
	return r.db.Delete(&entities.Job{}, id).Error
}

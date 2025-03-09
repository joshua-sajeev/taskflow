package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Predefined errors to reduce heap allocations
var ErrInvalidStatus = errors.New("invalid status")
var ErrEmptyTask = errors.New("task cannot be empty")

type Job struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Task      string    `gorm:"not null"`
	Status    string    `gorm:"default:pending"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// NewJob creates a new Job instance
func NewJob(task string) Job {
	return Job{
		Task:      task,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// UpdateStatus updates job status and validates input
func (j *Job) UpdateStatus(status string) error {
	if status != "pending" && status != "processing" && status != "completed" && status != "failed" {
		return ErrInvalidStatus
	}
	j.Status = status
	j.UpdatedAt = time.Now()
	return nil
}

func (j *Job) Execute() {
	j.UpdateStatus("processing")
	time.Sleep(2 * time.Second) // Simulating job execution
	j.UpdateStatus("completed")

}

// Migrate the Job table
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&Job{})
}

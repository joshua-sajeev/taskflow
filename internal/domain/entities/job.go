package entities

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Predefined errors to reduce heap allocations
var ErrInvalidStatus = errors.New("invalid status")
var ErrEmptyTask = errors.New("task cannot be empty")

type Job struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Task      string    `gorm:"not null"`
	Status    string    `gorm:"default:pending"` // pending, processing, completed, failed
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// NewJob returns a value type instead of a pointer to reduce heap escapes
func NewJob(task string) Job {
	return Job{
		Task:      task,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// UpdateStatus updates job status and validates input
func (j *Job) UpdateStatus(status string) (Job, error) {
	if status != "pending" && status != "processing" && status != "completed" && status != "failed" {
		return *j, ErrInvalidStatus
	}
	j.Status = status
	j.UpdatedAt = time.Now()
	return *j, nil
}

// Migrate without pointer to db to avoid heap escape
func Migrate(db gorm.DB) error {
	return db.AutoMigrate(&Job{})
}

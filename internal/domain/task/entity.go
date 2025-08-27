package task

import "time"

type Task struct {
	ID        int       `json:"id" binding:"gte=1" example:"1"`
	Task      string    `json:"task" binding:"required" example:"Buy milk"`
	Status    string    `json:"status" binding:"required" example:"pending"`
	CreatedAt time.Time `json:"created_at" example:"2025-08-27 10:35:16.263"`
}

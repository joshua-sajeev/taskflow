package task

import "time"

type Task struct {
	ID        int       `json:"id"`
	Task      string    `json:"task"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

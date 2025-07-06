package models

import "time"

type Status string

const (
	Pending   Status = "pending"
	Completed Status = "completed"
	Failed    Status = "failed"
)

type Task struct {
	ID          string        `json:"task_id"`
	Status      Status        `json:"status"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"created_at"`
	Duration    time.Duration `json:"duration"`
}

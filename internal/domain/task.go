package domain

import "time"

// TaskStatus represents the status of a task.
type TaskStatus string

const (
	StatusPending    TaskStatus = "PENDING"
	StatusInProgress TaskStatus = "IN_PROGRESS"
	StatusDone       TaskStatus = "DONE"
)

// Task represents a task entity in the domain.
type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Status      TaskStatus `json:"status"`
	DueDate     time.Time  `json:"due_date"`
}

// Validation error messages
const (
	ErrTitleRequired   = "title is required"
	ErrDueDateRequired = "due_date is required"
	ErrDueDatePast     = "due_date must be in the future"
	ErrStatusInvalid   = "invalid status"
)

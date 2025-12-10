package repository

import "github.com/gauravpandey771/task-api/internal/domain"

// TaskRepository interface (empty here, just for completeness)
// Real implementation is in task_repository_memory.go
type TaskRepository interface {
	domain.TaskRepository
}
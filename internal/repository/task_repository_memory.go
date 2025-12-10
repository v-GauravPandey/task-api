package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gauravpandey771/task-api/internal/domain"
	pkgerrors "github.com/gauravpandey771/task-api/pkg/errors"
)

// InMemoryTaskRepository is an in-memory implementation of TaskRepository.
type InMemoryTaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

// NewInMemoryTaskRepository creates a new in-memory repository.
func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[string]*domain.Task),
	}
}

// Create adds a new task to the repository.
func (r *InMemoryTaskRepository) Create(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Generate UUID for the task
	id := uuid.NewString()
	task.ID = id
	r.tasks[id] = task

	return nil
}

// GetByID retrieves a task by its ID.
func (r *InMemoryTaskRepository) GetByID(id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.tasks[id]
	if !ok {
		return nil, pkgerrors.NewNotFoundError("task not found")
	}

	// Return a copy to prevent external mutation
	copy := *task
	return &copy, nil
}

// Update updates an existing task.
func (r *InMemoryTaskRepository) Update(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[task.ID]; !ok {
		return pkgerrors.NewNotFoundError("task not found")
	}

	// Store a copy
	copy := *task
	r.tasks[task.ID] = &copy

	return nil
}

// Delete removes a task from the repository.
func (r *InMemoryTaskRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[id]; !ok {
		return pkgerrors.NewNotFoundError("task not found")
	}

	delete(r.tasks, id)
	return nil
}

// ListAll retrieves all tasks from the repository.
func (r *InMemoryTaskRepository) ListAll() ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]*domain.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		copy := *t
		out = append(out, &copy)
	}

	return out, nil
}
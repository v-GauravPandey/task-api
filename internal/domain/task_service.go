package domain

import (
	"sort"
	"time"

	pkgerrors "github.com/gauravpandey771/task-api/pkg/errors"
)

// TaskRepository defines the persistence interface.
type TaskRepository interface {
	Create(task *Task) error
	GetByID(id string) (*Task, error)
	Update(task *Task) error
	Delete(id string) error
	ListAll() ([]*Task, error)
}

// TaskService defines the business logic interface.
type TaskService interface {
	CreateTask(input CreateTaskInput) (*Task, error)
	GetTask(id string) (*Task, error)
	UpdateTask(id string, input UpdateTaskInput) (*Task, error)
	DeleteTask(id string) error
	ListTasks(filter TaskFilter) ([]*Task, error)
}

// CreateTaskInput is the input for creating a task.
type CreateTaskInput struct {
	Title       string
	Description string
	Status      *TaskStatus
	DueDate     time.Time
}

// UpdateTaskInput is the input for updating a task (all fields optional).
type UpdateTaskInput struct {
	Title       *string
	Description *string
	Status      *TaskStatus
	DueDate     *time.Time
}

// TaskFilter is used for listing tasks with filters and pagination.
type TaskFilter struct {
	Status   *TaskStatus
	Page     int
	PageSize int
}

// taskService implements TaskService interface.
type taskService struct {
	repo TaskRepository
}

// NewTaskService creates and returns a new TaskService.
func NewTaskService(repo TaskRepository) TaskService {
	return &taskService{repo: repo}
}

// CreateTask creates a new task with validation.
func (s *taskService) CreateTask(input CreateTaskInput) (*Task, error) {
	// Validate title
	if input.Title == "" {
		return nil, pkgerrors.NewValidationError(ErrTitleRequired)
	}

	// Validate due date
	if input.DueDate.IsZero() {
		return nil, pkgerrors.NewValidationError(ErrDueDateRequired)
	}
	if !input.DueDate.After(time.Now()) {
		return nil, pkgerrors.NewValidationError(ErrDueDatePast)
	}

	// Set default status or validate provided status
	status := StatusPending
	if input.Status != nil {
		if !isValidStatus(*input.Status) {
			return nil, pkgerrors.NewValidationError(ErrStatusInvalid)
		}
		status = *input.Status
	}

	// Create task entity
	task := &Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      status,
		DueDate:     input.DueDate,
	}

	// Persist
	if err := s.repo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

// GetTask retrieves a task by ID.
func (s *taskService) GetTask(id string) (*Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// UpdateTask updates an existing task with partial or full updates.
func (s *taskService) UpdateTask(id string, input UpdateTaskInput) (*Task, error) {
	// Get existing task
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update title
	if input.Title != nil {
		if *input.Title == "" {
			return nil, pkgerrors.NewValidationError(ErrTitleRequired)
		}
		task.Title = *input.Title
	}

	// Update description
	if input.Description != nil {
		task.Description = *input.Description
	}

	// Update status
	if input.Status != nil {
		if !isValidStatus(*input.Status) {
			return nil, pkgerrors.NewValidationError(ErrStatusInvalid)
		}
		task.Status = *input.Status
	}

	// Update due date
	if input.DueDate != nil {
		if input.DueDate.IsZero() {
			return nil, pkgerrors.NewValidationError(ErrDueDateRequired)
		}
		if !input.DueDate.After(time.Now()) {
			return nil, pkgerrors.NewValidationError(ErrDueDatePast)
		}
		task.DueDate = *input.DueDate
	}

	// Persist
	if err := s.repo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

// DeleteTask deletes a task by ID.
func (s *taskService) DeleteTask(id string) error {
	return s.repo.Delete(id)
}

// ListTasks lists all tasks with optional filtering and pagination.
func (s *taskService) ListTasks(filter TaskFilter) ([]*Task, error) {
	// Get all tasks
	tasks, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}

	// Filter by status if provided
	if filter.Status != nil {
		filtered := make([]*Task, 0, len(tasks))
		for _, t := range tasks {
			if t.Status == *filter.Status {
				filtered = append(filtered, t)
			}
		}
		tasks = filtered
	}

	// Sort by due date
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].DueDate.Before(tasks[j].DueDate)
	})

	// Apply pagination
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	size := filter.PageSize
	if size <= 0 {
		size = 10
	}

	start := (page - 1) * size
	if start >= len(tasks) {
		return []*Task{}, nil // Empty result if page out of range
	}

	end := start + size
	if end > len(tasks) {
		end = len(tasks)
	}

	return tasks[start:end], nil
}

// isValidStatus checks if a status string is valid.
func isValidStatus(s TaskStatus) bool {
	switch s {
	case StatusPending, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

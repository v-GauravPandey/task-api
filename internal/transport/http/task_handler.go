package http

import (
	"time"

	"github.com/gauravpandey771/task-api/internal/domain"
	pkgerrors "github.com/gauravpandey771/task-api/pkg/errors"
	"github.com/gofiber/fiber/v2"
)

// TaskHandler handles HTTP requests for tasks.
type TaskHandler struct {
	service domain.TaskService
}

// Request/Response DTOs
type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"due_date"` // ISO8601 format
}

type updateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	DueDate     *string `json:"due_date"` // ISO8601 format
}

// NewTaskHandler creates a new TaskHandler.
func NewTaskHandler(service domain.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// RegisterRoutes registers all task routes with a Fiber router.
func (h *TaskHandler) RegisterRoutes(r fiber.Router) {
	r.Post("/tasks", h.CreateTask)
	r.Get("/tasks/:id", h.GetTask)
	r.Put("/tasks/:id", h.UpdateTask)
	r.Delete("/tasks/:id", h.DeleteTask)
	r.Get("/tasks", h.ListTasks)
}

// CreateTask handles POST /tasks
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var req createTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}

	// Parse status if provided
	var statusPtr *domain.TaskStatus
	if req.Status != "" {
		s := domain.TaskStatus(req.Status)
		statusPtr = &s
	}

	// Parse due date if provided
	var due time.Time
	var err error
	if req.DueDate != "" {
		due, err = time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid due_date format, expected RFC3339")
		}
	}

	// Create task via service
	task, err := h.service.CreateTask(domain.CreateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      statusPtr,
		DueDate:     due,
	})
	if err != nil {
		if pkgerrors.IsValidation(err) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

// GetTask handles GET /tasks/:id
func (h *TaskHandler) GetTask(c *fiber.Ctx) error {
	id := c.Params("id")

	task, err := h.service.GetTask(id)
	if err != nil {
		if pkgerrors.IsNotFound(err) {
			return fiber.NewError(fiber.StatusNotFound, "task not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.JSON(task)
}

// UpdateTask handles PUT /tasks/:id
func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")

	var req updateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid JSON body")
	}

	// Parse status if provided
	var statusPtr *domain.TaskStatus
	if req.Status != nil {
		s := domain.TaskStatus(*req.Status)
		statusPtr = &s
	}

	// Parse due date if provided
	var duePtr *time.Time
	if req.DueDate != nil {
		if *req.DueDate != "" {
			d, err := time.Parse(time.RFC3339, *req.DueDate)
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "invalid due_date format, expected RFC3339")
			}
			duePtr = &d
		}
	}

	// Update task via service
	task, err := h.service.UpdateTask(id, domain.UpdateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      statusPtr,
		DueDate:     duePtr,
	})
	if err != nil {
		if pkgerrors.IsValidation(err) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if pkgerrors.IsNotFound(err) {
			return fiber.NewError(fiber.StatusNotFound, "task not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.JSON(task)
}

// DeleteTask handles DELETE /tasks/:id
func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.service.DeleteTask(id)
	if err != nil {
		if pkgerrors.IsNotFound(err) {
			return fiber.NewError(fiber.StatusNotFound, "task not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListTasks handles GET /tasks with optional filters
func (h *TaskHandler) ListTasks(c *fiber.Ctx) error {
	// Parse query parameters
	statusStr := c.Query("status")
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)

	// Parse status filter if provided
	var statusPtr *domain.TaskStatus
	if statusStr != "" {
		s := domain.TaskStatus(statusStr)
		statusPtr = &s
	}

	// List tasks via service
	tasks, err := h.service.ListTasks(domain.TaskFilter{
		Status:   statusPtr,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		if pkgerrors.IsValidation(err) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "internal error")
	}

	return c.JSON(tasks)
}

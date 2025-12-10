package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gauravpandey771/task-api/internal/domain"
	"github.com/gauravpandey771/task-api/internal/repository"
	httphandler "github.com/gauravpandey771/task-api/internal/transport/http"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Helper to create test app using Fiber
func newFiberTestApp() *fiber.App {
	repo := repository.NewInMemoryTaskRepository()
	svc := domain.NewTaskService(repo)
	handler := httphandler.NewTaskHandler(svc)
	return httphandler.NewApp(handler)
}

// TestCreateTaskValidation_EmptyTitle tests 400 error for empty title
func TestCreateTaskValidation_EmptyTitle(t *testing.T) {
	app := newFiberTestApp()
	due := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	body := map[string]any{"title": "", "due_date": due}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 5000)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// TestCreateTaskValidation_MissingDueDate tests 400 error for missing due date
func TestCreateTaskValidation_MissingDueDate(t *testing.T) {
	app := newFiberTestApp()
	body := map[string]any{"title": "Task"}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 5000)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// TestCreateTaskValidation_PastDueDate tests 400 error for past due date
func TestCreateTaskValidation_PastDueDate(t *testing.T) {
	app := newFiberTestApp()
	due := time.Now().Add(-24 * time.Hour).UTC().Format(time.RFC3339)
	body := map[string]any{"title": "Task", "due_date": due}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 5000)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// TestCreateTaskValidation_InvalidStatus tests 400 error for invalid status
func TestCreateTaskValidation_InvalidStatus(t *testing.T) {
	app := newFiberTestApp()
	due := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	body := map[string]any{"title": "Task", "status": "INVALID", "due_date": due}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 5000)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// TestCreateTask_WithAllFields tests creating task with all fields
func TestCreateTask_WithAllFields(t *testing.T) {
	app := newFiberTestApp()
	due := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	body := map[string]any{"title": "Complete", "description": "Full", "status": "IN_PROGRESS", "due_date": due}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 5000)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	respBody, _ := io.ReadAll(resp.Body)
	var created map[string]any
	json.Unmarshal(respBody, &created)
	assert.Equal(t, "Complete", created["title"])
	assert.Equal(t, "Full", created["description"])
	assert.Equal(t, "IN_PROGRESS", created["status"])
}

// TestError_GetTask_NotFound tests 404 for non-existent task (different from service test)
func TestError_GetTask_NotFound(t *testing.T) {
	app := newFiberTestApp()
	req, _ := http.NewRequest(http.MethodGet, "/tasks/non-existent", nil)
	resp, _ := app.Test(req, 5000)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// TestError_UpdateTask_NotFound tests 404 when updating non-existent task (different from service test)
func TestError_UpdateTask_NotFound(t *testing.T) {
	app := newFiberTestApp()
	body := map[string]any{"title": "Updated"}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, "/tasks/non-existent", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 5000)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// TestUpdateTask_InvalidStatus tests 400 error for invalid status in update
func TestUpdateTask_InvalidStatus(t *testing.T) {
	app := newFiberTestApp()
	due := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	body := map[string]any{"title": "Task", "due_date": due}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 5000)
	respBody, _ := io.ReadAll(resp.Body)
	var created map[string]any
	json.Unmarshal(respBody, &created)
	id := created["id"].(string)

	updateBody := map[string]any{"status": "INVALID"}
	updateB, _ := json.Marshal(updateBody)
	updateReq, _ := http.NewRequest(http.MethodPut, "/tasks/"+id, bytes.NewReader(updateB))
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp, _ := app.Test(updateReq, 5000)
	assert.Equal(t, http.StatusBadRequest, updateResp.StatusCode)
}

// TestError_DeleteTask_NotFound tests 404 when deleting non-existent task (different from service test)
func TestError_DeleteTask_NotFound(t *testing.T) {
	app := newFiberTestApp()
	req, _ := http.NewRequest(http.MethodDelete, "/tasks/non-existent", nil)
	resp, _ := app.Test(req, 5000)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// TestHandler_ListTasks_FilterByStatus tests filtering in list endpoint
func TestHandler_ListTasks_FilterByStatus(t *testing.T) {
	app := newFiberTestApp()
	due := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)

	body1 := map[string]any{"title": "Pending", "due_date": due}
	b1, _ := json.Marshal(body1)
	req1, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b1))
	req1.Header.Set("Content-Type", "application/json")
	app.Test(req1, 5000)

	body2 := map[string]any{"title": "Done", "status": "DONE", "due_date": due}
	b2, _ := json.Marshal(body2)
	req2, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b2))
	req2.Header.Set("Content-Type", "application/json")
	app.Test(req2, 5000)

	req, _ := http.NewRequest(http.MethodGet, "/tasks?status=DONE", nil)
	resp, _ := app.Test(req, 5000)
	respBody, _ := io.ReadAll(resp.Body)
	var tasks []map[string]any
	json.Unmarshal(respBody, &tasks)
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, "Done", tasks[0]["title"])
}

// TestHandler_ListTasks_Pagination tests pagination parameters
func TestHandler_ListTasks_Pagination(t *testing.T) {
	app := newFiberTestApp()
	due := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	for i := 0; i < 15; i++ {
		body := map[string]any{"title": "Task", "due_date": due}
		b, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		app.Test(req, 5000)
	}

	req1, _ := http.NewRequest(http.MethodGet, "/tasks?page=1&page_size=10", nil)
	resp1, _ := app.Test(req1, 5000)
	body1, _ := io.ReadAll(resp1.Body)
	var page1 []map[string]any
	json.Unmarshal(body1, &page1)
	assert.Equal(t, 10, len(page1))

	req2, _ := http.NewRequest(http.MethodGet, "/tasks?page=2&page_size=10", nil)
	resp2, _ := app.Test(req2, 5000)
	body2, _ := io.ReadAll(resp2.Body)
	var page2 []map[string]any
	json.Unmarshal(body2, &page2)
	assert.Equal(t, 5, len(page2))
}

// TestListTasks_DefaultValues tests default pagination values
func TestListTasks_DefaultValues(t *testing.T) {
	app := newFiberTestApp()
	due := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	for i := 0; i < 5; i++ {
		body := map[string]any{"title": "Task", "due_date": due}
		b, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		app.Test(req, 5000)
	}
	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	resp, _ := app.Test(req, 5000)
	respBody, _ := io.ReadAll(resp.Body)
	var tasks []map[string]any
	json.Unmarshal(respBody, &tasks)
	assert.Equal(t, 5, len(tasks))
}

// TestHandler_ListTasks_SortedByDueDate tests sorting by due date
func TestHandler_ListTasks_SortedByDueDate(t *testing.T) {
	app := newFiberTestApp()
	due2 := time.Now().Add(48 * time.Hour).UTC().Format(time.RFC3339)
	body2 := map[string]any{"title": "Task 2", "due_date": due2}
	b2, _ := json.Marshal(body2)
	req2, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b2))
	req2.Header.Set("Content-Type", "application/json")
	app.Test(req2, 5000)

	due1 := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	body1 := map[string]any{"title": "Task 1", "due_date": due1}
	b1, _ := json.Marshal(body1)
	req1, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b1))
	req1.Header.Set("Content-Type", "application/json")
	app.Test(req1, 5000)

	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	resp, _ := app.Test(req, 5000)
	respBody, _ := io.ReadAll(resp.Body)
	var tasks []map[string]any
	json.Unmarshal(respBody, &tasks)
	assert.Equal(t, "Task 1", tasks[0]["title"])
	assert.Equal(t, "Task 2", tasks[1]["title"])
}

// TestInvalidJSON tests invalid JSON handling
func TestInvalidJSON(t *testing.T) {
	app := newFiberTestApp()
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 5000)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

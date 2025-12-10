package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gauravpandey771/task-api/internal/domain"
	"github.com/gauravpandey771/task-api/internal/repository"
	httphandler "github.com/gauravpandey771/task-api/internal/transport/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper to create a test server
func newTestApp() *httptest.Server {
	repo := repository.NewInMemoryTaskRepository()
	svc := domain.NewTaskService(repo)
	handler := httphandler.NewTaskHandler(svc)
	app := httphandler.NewApp(handler)
	return httptest.NewServer(app)
}

// TestCreateAndGetTask tests end-to-end task creation and retrieval
func TestCreateAndGetTask(t *testing.T) {
	server := newTestApp()
	defer server.Close()

	due := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	body := map[string]any{
		"title":       "Integration Task",
		"description": "Test description",
		"due_date":    due,
	}
	b, _ := json.Marshal(body)

	// Create task
	resp, err := http.Post(server.URL+"/api/tasks", "application/json", bytes.NewReader(b))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))

	id, ok := created["id"].(string)
	require.True(t, ok)
	require.NotEmpty(t, id)

	// Get task
	getResp, err := http.Get(server.URL + "/api/tasks/" + id)
	require.NoError(t, err)
	defer getResp.Body.Close()

	assert.Equal(t, http.StatusOK, getResp.StatusCode)
	var got map[string]any
	require.NoError(t, json.NewDecoder(getResp.Body).Decode(&got))
	assert.Equal(t, "Integration Task", got["title"])
}

// TestListTasks tests listing all tasks
func TestListTasks(t *testing.T) {
	server := newTestApp()
	defer server.Close()

	// Create two tasks
	due := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	body := map[string]any{
		"title":    "Task 1",
		"due_date": due,
	}
	b, _ := json.Marshal(body)
	http.Post(server.URL+"/api/tasks", "application/json", bytes.NewReader(b))

	body["title"] = "Task 2"
	b, _ = json.Marshal(body)
	http.Post(server.URL+"/api/tasks", "application/json", bytes.NewReader(b))

	// List tasks
	resp, err := http.Get(server.URL + "/api/tasks")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var tasks []map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&tasks))
	assert.Equal(t, 2, len(tasks))
}

// TestUpdateTask tests task update
func TestUpdateTask(t *testing.T) {
	server := newTestApp()
	defer server.Close()

	// Create task
	due := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	body := map[string]any{
		"title":    "Original",
		"due_date": due,
	}
	b, _ := json.Marshal(body)

	resp, _ := http.Post(server.URL+"/api/tasks", "application/json", bytes.NewReader(b))
	var created map[string]any
	json.NewDecoder(resp.Body).Decode(&created)
	id := created["id"].(string)
	resp.Body.Close()

	// Update task
	updateBody := map[string]any{
		"title":  "Updated",
		"status": "IN_PROGRESS",
	}
	updateB, _ := json.Marshal(updateBody)

	req, _ := http.NewRequest(http.MethodPut, server.URL+"/api/tasks/"+id, bytes.NewReader(updateB))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	updateResp, _ := client.Do(req)
	defer updateResp.Body.Close()

	assert.Equal(t, http.StatusOK, updateResp.StatusCode)

	var updated map[string]any
	json.NewDecoder(updateResp.Body).Decode(&updated)
	assert.Equal(t, "Updated", updated["title"])
	assert.Equal(t, "IN_PROGRESS", updated["status"])
}

// TestDeleteTask tests task deletion
func TestDeleteTask(t *testing.T) {
	server := newTestApp()
	defer server.Close()

	// Create task
	due := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	body := map[string]any{
		"title":    "To Delete",
		"due_date": due,
	}
	b, _ := json.Marshal(body)

	resp, _ := http.Post(server.URL+"/api/tasks", "application/json", bytes.NewReader(b))
	var created map[string]any
	json.NewDecoder(resp.Body).Decode(&created)
	id := created["id"].(string)
	resp.Body.Close()

	// Delete task
	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/api/tasks/"+id, nil)
	client := &http.Client{}
	deleteResp, _ := client.Do(req)
	defer deleteResp.Body.Close()

	assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

	// Verify deletion
	getResp, _ := http.Get(server.URL + "/api/tasks/" + id)
	defer getResp.Body.Close()
	assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
}

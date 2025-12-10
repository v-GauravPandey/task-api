package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegration_CreateAndGetTask tests end-to-end task creation and retrieval
func TestIntegration_CreateAndGetTask(t *testing.T) {
	app := newFiberTestApp()
	due := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	body := map[string]any{"title": "Integration Task", "description": "Test description", "due_date": due}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 5000)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	var created map[string]any
	require.NoError(t, json.Unmarshal(respBody, &created))
	id, ok := created["id"].(string)
	require.True(t, ok)
	require.NotEmpty(t, id)

	getReq, _ := http.NewRequest(http.MethodGet, "/tasks/"+id, nil)
	getResp, _ := app.Test(getReq, 5000)
	assert.Equal(t, http.StatusOK, getResp.StatusCode)
	getRespBody, _ := io.ReadAll(getResp.Body)
	var got map[string]any
	require.NoError(t, json.Unmarshal(getRespBody, &got))
	assert.Equal(t, "Integration Task", got["title"])
}

// TestIntegration_ListTasks tests listing all tasks
func TestIntegration_ListTasks(t *testing.T) {
	app := newFiberTestApp()
	due := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)

	body := map[string]any{"title": "Task 1", "due_date": due}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	app.Test(req, 5000)

	body["title"] = "Task 2"
	b, _ = json.Marshal(body)
	req2, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b))
	req2.Header.Set("Content-Type", "application/json")
	app.Test(req2, 5000)

	listReq, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	resp, _ := app.Test(listReq, 5000)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	respBody, _ := io.ReadAll(resp.Body)
	var tasks []map[string]any
	require.NoError(t, json.Unmarshal(respBody, &tasks))
	assert.Equal(t, 2, len(tasks))
}

// TestIntegration_UpdateTask tests task update
func TestIntegration_UpdateTask(t *testing.T) {
	app := newFiberTestApp()
	due := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	body := map[string]any{"title": "Original", "due_date": due}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 5000)
	respBody, _ := io.ReadAll(resp.Body)
	var created map[string]any
	json.Unmarshal(respBody, &created)
	id := created["id"].(string)

	updateBody := map[string]any{"title": "Updated", "status": "IN_PROGRESS"}
	updateB, _ := json.Marshal(updateBody)
	updateReq, _ := http.NewRequest(http.MethodPut, "/tasks/"+id, bytes.NewReader(updateB))
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp, _ := app.Test(updateReq, 5000)
	assert.Equal(t, http.StatusOK, updateResp.StatusCode)

	updateRespBody, _ := io.ReadAll(updateResp.Body)
	var updated map[string]any
	json.Unmarshal(updateRespBody, &updated)
	assert.Equal(t, "Updated", updated["title"])
	assert.Equal(t, "IN_PROGRESS", updated["status"])
}

// TestIntegration_DeleteTask tests task deletion
func TestIntegration_DeleteTask(t *testing.T) {
	app := newFiberTestApp()
	due := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	body := map[string]any{"title": "To Delete", "due_date": due}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, 5000)
	respBody, _ := io.ReadAll(resp.Body)
	var created map[string]any
	json.Unmarshal(respBody, &created)
	id := created["id"].(string)

	delReq, _ := http.NewRequest(http.MethodDelete, "/tasks/"+id, nil)
	delResp, _ := app.Test(delReq, 5000)
	assert.Equal(t, http.StatusNoContent, delResp.StatusCode)

	getReq, _ := http.NewRequest(http.MethodGet, "/tasks/"+id, nil)
	getResp, _ := app.Test(getReq, 5000)
	assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
}

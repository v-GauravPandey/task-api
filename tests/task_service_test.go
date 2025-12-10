package tests

import (
	"testing"
	"time"

	"github.com/gauravpandey771/task-api/internal/domain"
	"github.com/gauravpandey771/task-api/internal/repository"
	pkgerrors "github.com/gauravpandey771/task-api/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper to create a test service
func newTestService() domain.TaskService {
	repo := repository.NewInMemoryTaskRepository()
	return domain.NewTaskService(repo)
}

// TestCreateTask_Success tests successful task creation
func TestCreateTask_Success(t *testing.T) {
	svc := newTestService()
	due := time.Now().Add(24 * time.Hour)

	task, err := svc.CreateTask(domain.CreateTaskInput{
		Title:   "Test Task",
		DueDate: due,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, task.ID)
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, domain.StatusPending, task.Status)
}

// TestCreateTask_WithDescription tests task creation with description
func TestCreateTask_WithDescription(t *testing.T) {
	svc := newTestService()
	due := time.Now().Add(24 * time.Hour)

	task, err := svc.CreateTask(domain.CreateTaskInput{
		Title:       "Task with desc",
		Description: "This is a description",
		DueDate:     due,
	})

	require.NoError(t, err)
	assert.Equal(t, "This is a description", task.Description)
}

// TestCreateTask_WithCustomStatus tests task creation with custom status
func TestCreateTask_WithCustomStatus(t *testing.T) {
	svc := newTestService()
	due := time.Now().Add(24 * time.Hour)
	status := domain.StatusInProgress

	task, err := svc.CreateTask(domain.CreateTaskInput{
		Title:   "In Progress Task",
		Status:  &status,
		DueDate: due,
	})

	require.NoError(t, err)
	assert.Equal(t, domain.StatusInProgress, task.Status)
}

// TestCreateTask_MissingTitle tests validation for missing title
func TestCreateTask_MissingTitle(t *testing.T) {
	svc := newTestService()
	due := time.Now().Add(24 * time.Hour)

	_, err := svc.CreateTask(domain.CreateTaskInput{
		Title:   "",
		DueDate: due,
	})

	require.Error(t, err)
	assert.True(t, pkgerrors.IsValidation(err))
	assert.Equal(t, domain.ErrTitleRequired, err.Error())
}

// TestCreateTask_MissingDueDate tests validation for missing due date
func TestCreateTask_MissingDueDate(t *testing.T) {
	svc := newTestService()

	_, err := svc.CreateTask(domain.CreateTaskInput{
		Title:   "Task without date",
		DueDate: time.Time{},
	})

	require.Error(t, err)
	assert.True(t, pkgerrors.IsValidation(err))
	assert.Equal(t, domain.ErrDueDateRequired, err.Error())
}

// TestCreateTask_PastDueDate tests validation for past due date
func TestCreateTask_PastDueDate(t *testing.T) {
	svc := newTestService()
	due := time.Now().Add(-24 * time.Hour)

	_, err := svc.CreateTask(domain.CreateTaskInput{
		Title:   "Task",
		DueDate: due,
	})

	require.Error(t, err)
	assert.True(t, pkgerrors.IsValidation(err))
	assert.Equal(t, domain.ErrDueDatePast, err.Error())
}

// TestCreateTask_InvalidStatus tests validation for invalid status
func TestCreateTask_InvalidStatus(t *testing.T) {
	svc := newTestService()
	due := time.Now().Add(24 * time.Hour)
	invalidStatus := domain.TaskStatus("INVALID")

	_, err := svc.CreateTask(domain.CreateTaskInput{
		Title:   "Task",
		Status:  &invalidStatus,
		DueDate: due,
	})

	require.Error(t, err)
	assert.True(t, pkgerrors.IsValidation(err))
	assert.Equal(t, domain.ErrStatusInvalid, err.Error())
}

// TestGetTask_Success tests successful task retrieval
func TestGetTask_Success(t *testing.T) {
	svc := newTestService()
	due := time.Now().Add(24 * time.Hour)

	created, _ := svc.CreateTask(domain.CreateTaskInput{
		Title:   "Task to Get",
		DueDate: due,
	})

	got, err := svc.GetTask(created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
	assert.Equal(t, "Task to Get", got.Title)
}

// TestGetTask_NotFound tests retrieval of non-existent task
func TestGetTask_NotFound(t *testing.T) {
	svc := newTestService()

	_, err := svc.GetTask("non-existent")
	require.Error(t, err)
	assert.True(t, pkgerrors.IsNotFound(err))
}

// TestUpdateTask_Success tests successful task update
func TestUpdateTask_Success(t *testing.T) {
	svc := newTestService()
	due := time.Now().Add(24 * time.Hour)

	created, _ := svc.CreateTask(domain.CreateTaskInput{
		Title:   "Original Title",
		DueDate: due,
	})

	newTitle := "Updated Title"
	updated, err := svc.UpdateTask(created.ID, domain.UpdateTaskInput{
		Title: &newTitle,
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated Title", updated.Title)
}

// TestUpdateTask_NotFound tests update of non-existent task
func TestUpdateTask_NotFound(t *testing.T) {
	svc := newTestService()
	title := "Updated"

	_, err := svc.UpdateTask("non-existent", domain.UpdateTaskInput{
		Title: &title,
	})

	require.Error(t, err)
	assert.True(t, pkgerrors.IsNotFound(err))
}

// TestDeleteTask_Success tests successful task deletion
func TestDeleteTask_Success(t *testing.T) {
	svc := newTestService()
	due := time.Now().Add(24 * time.Hour)

	created, _ := svc.CreateTask(domain.CreateTaskInput{
		Title:   "Task to Delete",
		DueDate: due,
	})

	err := svc.DeleteTask(created.ID)
	require.NoError(t, err)

	// Verify it's deleted
	_, err = svc.GetTask(created.ID)
	assert.True(t, pkgerrors.IsNotFound(err))
}

// TestDeleteTask_NotFound tests deletion of non-existent task
func TestDeleteTask_NotFound(t *testing.T) {
	svc := newTestService()

	err := svc.DeleteTask("non-existent")
	require.Error(t, err)
	assert.True(t, pkgerrors.IsNotFound(err))
}

// TestListTasks_Empty tests listing with no tasks
func TestListTasks_Empty(t *testing.T) {
	svc := newTestService()

	tasks, err := svc.ListTasks(domain.TaskFilter{})
	require.NoError(t, err)
	assert.Equal(t, 0, len(tasks))
}

// TestListTasks_Multiple tests listing multiple tasks
func TestListTasks_Multiple(t *testing.T) {
	svc := newTestService()
	due1 := time.Now().Add(24 * time.Hour)
	due2 := time.Now().Add(48 * time.Hour)

	svc.CreateTask(domain.CreateTaskInput{Title: "Task 1", DueDate: due1})
	svc.CreateTask(domain.CreateTaskInput{Title: "Task 2", DueDate: due2})

	tasks, err := svc.ListTasks(domain.TaskFilter{})
	require.NoError(t, err)
	assert.Equal(t, 2, len(tasks))
}

// TestListTasks_SortedByDueDate tests that tasks are sorted by due date
func TestListTasks_SortedByDueDate(t *testing.T) {
	svc := newTestService()
	due2 := time.Now().Add(48 * time.Hour)
	due1 := time.Now().Add(24 * time.Hour)

	svc.CreateTask(domain.CreateTaskInput{Title: "Task Later", DueDate: due2})
	svc.CreateTask(domain.CreateTaskInput{Title: "Task Earlier", DueDate: due1})

	tasks, err := svc.ListTasks(domain.TaskFilter{})
	require.NoError(t, err)
	assert.Equal(t, "Task Earlier", tasks[0].Title)
	assert.Equal(t, "Task Later", tasks[1].Title)
}

// TestListTasks_FilterByStatus tests filtering by status
func TestListTasks_FilterByStatus(t *testing.T) {
	svc := newTestService()
	due := time.Now().Add(24 * time.Hour)
	status := domain.StatusDone

	svc.CreateTask(domain.CreateTaskInput{Title: "Pending Task", DueDate: due})
	svc.CreateTask(domain.CreateTaskInput{Title: "Done Task", Status: &status, DueDate: due})

	tasks, err := svc.ListTasks(domain.TaskFilter{Status: &status})
	require.NoError(t, err)
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, domain.StatusDone, tasks[0].Status)
}

// TestListTasks_Pagination tests pagination
func TestListTasks_Pagination(t *testing.T) {
	svc := newTestService()
	due := time.Now().Add(24 * time.Hour)

	for i := 0; i < 25; i++ {
		svc.CreateTask(domain.CreateTaskInput{Title: "Task", DueDate: due})
	}

	// First page
	tasks1, err := svc.ListTasks(domain.TaskFilter{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, 10, len(tasks1))

	// Second page
	tasks2, err := svc.ListTasks(domain.TaskFilter{Page: 2, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, 10, len(tasks2))

	// Third page (partial)
	tasks3, err := svc.ListTasks(domain.TaskFilter{Page: 3, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, 5, len(tasks3))
}

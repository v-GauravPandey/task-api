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

// TestRepository_CreateTask tests task creation in repository
func TestRepository_CreateTask(t *testing.T) {
	repo := repository.NewInMemoryTaskRepository()
	task := &domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      domain.StatusPending,
		DueDate:     time.Now().Add(24 * time.Hour),
	}

	err := repo.Create(task)
	require.NoError(t, err)
	assert.NotEmpty(t, task.ID)
}

// TestRepository_GetByID tests retrieving a task by ID
func TestRepository_GetByID(t *testing.T) {
	repo := repository.NewInMemoryTaskRepository()
	task := &domain.Task{
		Title:   "Test Task",
		Status:  domain.StatusPending,
		DueDate: time.Now().Add(24 * time.Hour),
	}

	repo.Create(task)
	retrieved, err := repo.GetByID(task.ID)

	require.NoError(t, err)
	assert.Equal(t, task.Title, retrieved.Title)
	assert.Equal(t, task.ID, retrieved.ID)
}

// TestRepository_GetByID_NotFound tests retrieval of non-existent task
func TestRepository_GetByID_NotFound(t *testing.T) {
	repo := repository.NewInMemoryTaskRepository()

	_, err := repo.GetByID("non-existent-id")
	require.Error(t, err)
	assert.True(t, pkgerrors.IsNotFound(err))
}

// TestRepository_Update tests task update
func TestRepository_Update(t *testing.T) {
	repo := repository.NewInMemoryTaskRepository()
	task := &domain.Task{
		Title:   "Original Title",
		Status:  domain.StatusPending,
		DueDate: time.Now().Add(24 * time.Hour),
	}

	repo.Create(task)
	task.Title = "Updated Title"
	err := repo.Update(task)

	require.NoError(t, err)

	retrieved, _ := repo.GetByID(task.ID)
	assert.Equal(t, "Updated Title", retrieved.Title)
}

// TestRepository_Update_NotFound tests update of non-existent task
func TestRepository_Update_NotFound(t *testing.T) {
	repo := repository.NewInMemoryTaskRepository()
	task := &domain.Task{
		ID:      "non-existent",
		Title:   "Test",
		Status:  domain.StatusPending,
		DueDate: time.Now().Add(24 * time.Hour),
	}

	err := repo.Update(task)
	require.Error(t, err)
	assert.True(t, pkgerrors.IsNotFound(err))
}

// TestRepository_Delete tests task deletion
func TestRepository_Delete(t *testing.T) {
	repo := repository.NewInMemoryTaskRepository()
	task := &domain.Task{
		Title:   "Task to Delete",
		Status:  domain.StatusPending,
		DueDate: time.Now().Add(24 * time.Hour),
	}

	repo.Create(task)
	err := repo.Delete(task.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(task.ID)
	assert.True(t, pkgerrors.IsNotFound(err))
}

// TestRepository_Delete_NotFound tests deletion of non-existent task
func TestRepository_Delete_NotFound(t *testing.T) {
	repo := repository.NewInMemoryTaskRepository()

	err := repo.Delete("non-existent-id")
	require.Error(t, err)
	assert.True(t, pkgerrors.IsNotFound(err))
}

// TestRepository_ListAll tests listing all tasks
func TestRepository_ListAll(t *testing.T) {
	repo := repository.NewInMemoryTaskRepository()

	task1 := &domain.Task{
		Title:   "Task 1",
		Status:  domain.StatusPending,
		DueDate: time.Now().Add(24 * time.Hour),
	}
	task2 := &domain.Task{
		Title:   "Task 2",
		Status:  domain.StatusInProgress,
		DueDate: time.Now().Add(48 * time.Hour),
	}

	repo.Create(task1)
	repo.Create(task2)

	tasks, err := repo.ListAll()
	require.NoError(t, err)
	assert.Equal(t, 2, len(tasks))
}

// TestRepository_ListAll_Empty tests listing with no tasks
func TestRepository_ListAll_Empty(t *testing.T) {
	repo := repository.NewInMemoryTaskRepository()

	tasks, err := repo.ListAll()
	require.NoError(t, err)
	assert.Equal(t, 0, len(tasks))
}

// TestRepository_Isolation tests that repository stores references properly
func TestRepository_Isolation(t *testing.T) {
	repo := repository.NewInMemoryTaskRepository()
	task := &domain.Task{
		Title:   "Original",
		Status:  domain.StatusPending,
		DueDate: time.Now().Add(24 * time.Hour),
	}

	repo.Create(task)
	taskID := task.ID

	// Verify the task was created with correct values
	retrieved, _ := repo.GetByID(taskID)
	assert.Equal(t, "Original", retrieved.Title)

	// Update the original task struct and verify GetByID returns a copy
	task.Title = "Modified"

	// GetByID returns a shallow copy of the stored value,
	// so it should show the modified title since we modified the original pointer
	// This is expected behavior as we store pointers
	retrieved2, _ := repo.GetByID(taskID)
	assert.Equal(t, "Modified", retrieved2.Title)

	// But if we modify the retrieved copy, it shouldn't affect the next retrieval
	retrieved2.Title = "AnotherModification"
	retrieved3, _ := repo.GetByID(taskID)
	// Since GetByID makes a shallow copy, this should still be "Modified"
	assert.Equal(t, "Modified", retrieved3.Title)
}

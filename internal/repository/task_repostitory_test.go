package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Andrew1996-la/timo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	require.NoError(t, storage.Migrate(db))

	return db
}

func TestTaskRepository_Create_Success(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	repo := NewTaskRepository(db)

	task, err := repo.Create(ctx, "Learn repository tests")

	require.NoError(t, err)
	require.NotNil(t, task)

	assert.NotZero(t, task.Id)
	assert.Equal(t, "Learn repository tests", task.Title)
	assert.Equal(t, 0, task.SpentSeconds)
	assert.Nil(t, task.DeletedAt)
	assert.False(t, task.CreatedAt.IsZero())
}

func TestTaskRepository_GetByID_Success(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	repo := NewTaskRepository(db)

	created, err := repo.Create(ctx, "Test task")
	require.NoError(t, err)

	task, err := repo.GetByID(ctx, created.Id)

	require.NoError(t, err)
	require.NotNil(t, task)

	assert.Equal(t, created.Id, task.Id)
	assert.Equal(t, created.Title, task.Title)
	assert.Equal(t, created.SpentSeconds, task.SpentSeconds)
	assert.Equal(t, created.CreatedAt.Unix(), task.CreatedAt.Unix())
	assert.Equal(t, created.DeletedAt, task.DeletedAt)
}

func TestTaskRepository_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	repo := NewTaskRepository(db)

	task, err := repo.GetByID(ctx, 999)

	require.Error(t, err)
	assert.Nil(t, task)
	assert.True(t, errors.Is(err, ErrTaskNotFound))
}

func TestTaskRepository_GetAll_Empty(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	repo := NewTaskRepository(db)

	tasks, err := repo.GetAll(ctx)

	require.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestTaskRepository_GetAll_Success(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	repo := NewTaskRepository(db)

	task1, err := repo.Create(ctx, "First task")
	require.NoError(t, err)

	task2, err := repo.Create(ctx, "Second task")
	require.NoError(t, err)

	task3, err := repo.Create(ctx, "Third task")
	require.NoError(t, err)

	tasks, err := repo.GetAll(ctx)

	require.NoError(t, err)
	assert.Len(t, tasks, 3)

	assert.Equal(t, task3.Id, tasks[0].Id)
	assert.Equal(t, task3.Title, tasks[0].Title)

	assert.Equal(t, task2.Id, tasks[1].Id)
	assert.Equal(t, task2.Title, tasks[1].Title)

	assert.Equal(t, task1.Id, tasks[2].Id)
	assert.Equal(t, task1.Title, tasks[2].Title)

	for _, task := range tasks {
		assert.NotZero(t, task.Id)
		assert.NotEmpty(t, task.Title)
		assert.False(t, task.CreatedAt.IsZero())
		assert.Nil(t, task.DeletedAt)
		assert.Equal(t, 0, task.SpentSeconds)
	}
}

func TestTaskRepository_Delete_Success(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	repo := NewTaskRepository(db)

	task, err := repo.Create(ctx, "Task to delete")
	require.NoError(t, err)

	err = repo.Delete(ctx, task.Id)
	require.NoError(t, err)

	deletedTask, err := repo.GetByID(ctx, task.Id)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrTaskNotFound))
	assert.Nil(t, deletedTask)

	tasks, err := repo.GetAll(ctx)
	require.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestTaskRepository_Delete_NotFound(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	repo := NewTaskRepository(db)

	err := repo.Delete(ctx, 999)

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrTaskNotFound))
}

func TestTaskRepository_Delete_AlreadyDeleted(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	repo := NewTaskRepository(db)

	task, err := repo.Create(ctx, "Task to delete twice")
	require.NoError(t, err)

	err = repo.Delete(ctx, task.Id)
	require.NoError(t, err)

	err = repo.Delete(ctx, task.Id)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrTaskNotFound))
}

func TestTaskRepository_AddTime_Success(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	repo := NewTaskRepository(db)

	task, err := repo.Create(ctx, "Task with time")
	require.NoError(t, err)

	err = repo.AddTime(ctx, task.Id, 60)
	require.NoError(t, err)

	task, err = repo.GetByID(ctx, task.Id)
	require.NoError(t, err)
	assert.Equal(t, 60, task.SpentSeconds)

	err = repo.AddTime(ctx, task.Id, 30)
	require.NoError(t, err)

	task, err = repo.GetByID(ctx, task.Id)
	require.NoError(t, err)
	assert.Equal(t, 90, task.SpentSeconds)
}

func TestTaskRepository_AddTime_NotFound(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	repo := NewTaskRepository(db)

	err := repo.AddTime(ctx, 999, 60)

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrTaskNotFound))
}

func TestTaskRepository_AddTime_ToDeletedTask(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	repo := NewTaskRepository(db)

	task, err := repo.Create(ctx, "Task to delete")
	require.NoError(t, err)

	err = repo.Delete(ctx, task.Id)
	require.NoError(t, err)

	err = repo.AddTime(ctx, task.Id, 60)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrTaskNotFound))
}

func TestTaskRepository_AddTime_Negative(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	repo := NewTaskRepository(db)

	task, err := repo.Create(ctx, "Task with negative time")
	require.NoError(t, err)

	err = repo.AddTime(ctx, task.Id, -10)
	require.NoError(t, err)

	task, err = repo.GetByID(ctx, task.Id)
	require.NoError(t, err)
	assert.Equal(t, -10, task.SpentSeconds)
}

func TestTaskRepository_Integration_Workflow(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	repo := NewTaskRepository(db)

	task1, err := repo.Create(ctx, "Task 1")
	require.NoError(t, err)

	task2, err := repo.Create(ctx, "Task 2")
	require.NoError(t, err)

	err = repo.AddTime(ctx, task1.Id, 120)
	require.NoError(t, err)

	task1, err = repo.GetByID(ctx, task1.Id)
	require.NoError(t, err)
	assert.Equal(t, 120, task1.SpentSeconds)

	err = repo.Delete(ctx, task2.Id)
	require.NoError(t, err)

	tasks, err := repo.GetAll(ctx)
	require.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, task1.Id, tasks[0].Id)

	_, err = repo.GetByID(ctx, task2.Id)
	assert.True(t, errors.Is(err, ErrTaskNotFound))

	err = repo.AddTime(ctx, task1.Id, 30)
	require.NoError(t, err)

	task1, err = repo.GetByID(ctx, task1.Id)
	require.NoError(t, err)
	assert.Equal(t, 150, task1.SpentSeconds)

	err = repo.Delete(ctx, task1.Id)
	require.NoError(t, err)

	tasks, err = repo.GetAll(ctx)
	require.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestTaskRepository_Concurrent_Create(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	repo := NewTaskRepository(db)

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			task, err := repo.Create(ctx, "Concurrent task")
			assert.NoError(t, err)
			assert.NotZero(t, task.Id)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	tasks, err := repo.GetAll(ctx)
	require.NoError(t, err)
	assert.Len(t, tasks, 10)
}

func TestTaskRepository_Concurrent_AddTime(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	repo := NewTaskRepository(db)

	task, err := repo.Create(ctx, "Concurrent time task")
	require.NoError(t, err)

	done := make(chan bool)
	for i := 0; i < 5; i++ {
		go func() {
			err := repo.AddTime(ctx, task.Id, 10)
			assert.NoError(t, err)
			done <- true
		}()
	}

	for i := 0; i < 5; i++ {
		<-done
	}

	task, err = repo.GetByID(ctx, task.Id)
	require.NoError(t, err)
	assert.Equal(t, 50, task.SpentSeconds)
}

func TestTaskRepository_DeletedAt_Formatting(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	repo := NewTaskRepository(db)

	task, err := repo.Create(ctx, "Task with deletion time")
	require.NoError(t, err)

	err = repo.Delete(ctx, task.Id)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, task.Id)
	assert.True(t, errors.Is(err, ErrTaskNotFound))

	tasks, err := repo.GetAll(ctx)
	require.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestTaskRepository_Timestamp_Consistency(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	repo := NewTaskRepository(db)

	beforeCreate := time.Now()
	task, err := repo.Create(ctx, "Timestamp test")
	afterCreate := time.Now()
	require.NoError(t, err)

	assert.True(t, task.CreatedAt.After(beforeCreate) || task.CreatedAt.Equal(beforeCreate))
	assert.True(t, task.CreatedAt.Before(afterCreate) || task.CreatedAt.Equal(afterCreate))

	beforeAddTime := task.CreatedAt
	err = repo.AddTime(ctx, task.Id, 30)
	require.NoError(t, err)

	task, err = repo.GetByID(ctx, task.Id)
	require.NoError(t, err)
	assert.Equal(t, beforeAddTime.Unix(), task.CreatedAt.Unix())
}
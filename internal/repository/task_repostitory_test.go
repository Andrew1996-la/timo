package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

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
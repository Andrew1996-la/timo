package storage

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "modernc.org/sqlite"
)

func setupMigrationTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	return db
}

func TestMigrate_CreatesTasksTable(t *testing.T) {
	db := setupMigrationTestDB(t)

	err := Migrate(db)

	require.NoError(t, err)

	var tableName string
	err = db.QueryRow(`
		SELECT name
		FROM sqlite_master
		WHERE type = 'table' AND name = 'tasks'
	`).Scan(&tableName)

	require.NoError(t, err)
	assert.Equal(t, "tasks", tableName)
}

func TestMigrate_CreatesMigrationsTable(t *testing.T) {
	db := setupMigrationTestDB(t)

	err := Migrate(db)

	require.NoError(t, err)

	var tableName string
	err = db.QueryRow(`
		SELECT name
		FROM sqlite_master
		WHERE type = 'table' AND name = 'migrations'
	`).Scan(&tableName)

	require.NoError(t, err)
	assert.Equal(t, "migrations", tableName)
}

func TestMigrate_SavesAppliedMigration(t *testing.T) {
	db := setupMigrationTestDB(t)

	err := Migrate(db)

	require.NoError(t, err)

	var migrationName string
	err = db.QueryRow(`
		SELECT name
		FROM migrations
		WHERE name = ?
	`, "001_create_tasks").Scan(&migrationName)

	require.NoError(t, err)
	assert.Equal(t, "001_create_tasks", migrationName)
}

func TestMigrate_CanRunTwice(t *testing.T) {
	db := setupMigrationTestDB(t)

	require.NoError(t, Migrate(db))
	require.NoError(t, Migrate(db))

	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM migrations
		WHERE name = ?
	`, "001_create_tasks").Scan(&count)

	require.NoError(t, err)
	assert.Equal(t, 1, count)
}
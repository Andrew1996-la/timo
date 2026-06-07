package storage

import (
	"database/sql"
	"errors"
	"fmt"
)

func Migrate(db *sql.DB) error {
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	for _, migration := range migrations {
		applied, err := isMigrationApplied(db, migration.Name)
		if err != nil {
			return fmt.Errorf(
				"check migration %s: %w",
				migration.Name,
				err,
			)
		}

		if applied {
			continue
		}

		if err := applyMigration(db, migration); err != nil {
			return err
		}
	}

	return nil
}

func createMigrationsTable(db *sql.DB) error {
	const query = `
		CREATE TABLE IF NOT EXISTS migrations (
			name TEXT PRIMARY KEY,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`

	_, err := db.Exec(query)

	return err
}

func isMigrationApplied(db *sql.DB, name string) (bool, error) {
	const query = `
		SELECT name
		FROM migrations
		WHERE name = ?
	`

	var migrationName string

	err := db.QueryRow(query, name).Scan(&migrationName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func applyMigration(db *sql.DB, migration Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf(
			"begin migration %s: %w",
			migration.Name,
			err,
		)
	}

	defer func() {
		_ = tx.Rollback()
	}()
	
	if _, err := tx.Exec(migration.Query); err != nil {
		return fmt.Errorf(
			"execute migration %s: %w",
			migration.Name,
			err,
		)
	}

	const insertQuery = `
		INSERT INTO migrations(name)
		VALUES(?)
	`

	if _, err := tx.Exec(insertQuery, migration.Name); err != nil {
		return fmt.Errorf(
			"save migration %s: %w",
			migration.Name,
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(
			"commit migration %s: %w",
			migration.Name,
			err,
		)
	}

	return nil
}

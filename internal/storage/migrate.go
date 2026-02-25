package storage

import (
	"database/sql"
	"fmt"
)

func Migrate(db *sql.DB) error {
	// создаём таблицу миграций
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS migrations (
		name TEXT PRIMARY KEY,
		applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		return err
	}

	for _, m := range migrations {
		var exists string
		err := db.QueryRow(`SELECT name FROM migrations WHERE name = ?`, m.Name).Scan(&exists)
		if err == sql.ErrNoRows {
			// миграция ещё не применена
			_, err := db.Exec(m.Query)
			if err != nil {
				return fmt.Errorf("failed migration %s: %w", m.Name, err)
			}
			// помечаем как применённую
			_, err = db.Exec(`INSERT INTO migrations(name) VALUES(?)`, m.Name)
			if err != nil {
				return fmt.Errorf("failed to insert migration record %s: %w", m.Name, err)
			}
			// fmt.Println("applied migration:", m.Name)
		} else if err != nil {
			return err
		} else {
			// миграция уже есть — пропускаем
			// fmt.Println("skipped migration:", m.Name)
		}
	}

	return nil
}

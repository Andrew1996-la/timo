package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const (
	appDirName = "timo"
	dbFileName = "timo.db"
	sqliteDriverName = "sqlite"
)

func NewSQLite() (*sql.DB, error) {
	dbPath, err := sqlitePath()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(sqliteDriverName, dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	if err := Migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate sqlite: %w", err)
	}

	return db, nil
}

func sqlitePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config dir: %w", err)
	}

	appDir := filepath.Join(configDir, appDirName)

	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", fmt.Errorf("create app dir: %w", err)
	}

	return filepath.Join(appDir, dbFileName), nil
}
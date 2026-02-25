package storage

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func NewSQLite() (*sql.DB, error) {
	// Получаем системную директорию для данных
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	// Создаём папку timo
	appDir := filepath.Join(dir, "timo")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, err
	}

	// Полный путь к БД
	dbPath := filepath.Join(appDir, "timo.db")

	// Открываем SQLite
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// применяем миграции
	if err := Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Andrew1996-la/timo/internal/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, title string) (*models.Task, error) {
	query := `
		INSERT INTO tasks (title, created_at)
		VALUES (?, CURRENT_TIMESTAMP)
	`

	res, err := r.db.Exec(query, title)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, int(id))
}

func (r *TaskRepository) Delete(ctx context.Context, id int) error {
	query := `
		UPDATE tasks
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("task not found or already deleted")
	}

	return nil
}

func (r *TaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	query := `
		SELECT id, title, created_at, deleted_at, spent_seconds
		FROM tasks
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var t models.Task
		if err := rows.Scan(
			&t.Id,
			&t.Title,
			&t.CreatedAt,
			&t.DeletedAt,
			&t.SpentSeconds,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id int) (*models.Task, error) {
	query := `
		SELECT id, title, created_at, deleted_at, spent_seconds
		FROM tasks
		WHERE id = ? AND deleted_at IS NULL
	`

	var t models.Task

	err := r.db.QueryRow(query, id).Scan(
		&t.Id,
		&t.Title,
		&t.CreatedAt,
		&t.DeletedAt,
		&t.SpentSeconds,
	)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TaskRepository) AddTime(
	ctx context.Context,
	id int,
	seconds int,
) error {
	query := `
		UPDATE tasks
		SET spent_seconds = spent_seconds + ?
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.Exec(query, seconds, id)
	return err
}

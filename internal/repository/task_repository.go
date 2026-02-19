package repository

import (
	"context"

	"github.com/Andrew1996-la/timo/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (t *TaskRepository) Create(ctx context.Context, title string) (*models.Task, error) {
	var task models.Task

	query := `
		INSERT INTO tasks (title)
		VALUES ($1)
		RETURNING id, title, created_at, deleted_at
	`

	err := t.db.QueryRow(ctx, query, title).
		Scan(&task.Id, &task.Title, &task.CreatedAt, &task.DeletedAt)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

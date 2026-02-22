package repository

import (
	"context"
	"errors"

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

func (t *TaskRepository) Delete(ctx context.Context, id int) error {
	query := `
		UPDATE tasks
		SET deleted_at = NOW()
		where id = $1 and deleted_at IS NULL
   `

	cmd, err := t.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("task not found or already deleted")
	}

	return nil
}

func (t *TaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	query := `
		SELECT id, title, created_at, deleted_at, spent_seconds
		FROM tasks
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := t.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task

		err = rows.Scan(&task.Id, &task.Title, &task.CreatedAt, &task.DeletedAt, &task.SpentSeconds)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (t *TaskRepository) GetByID(ctx context.Context, id int) (*models.Task, error) {
	query := `
		SELECT id, title, created_at, deleted_at, spent_seconds
		FROM tasks
		WHERE id = $1 AND deleted_at IS NULL
	`

	row := t.db.QueryRow(ctx, query, id)

	var task models.Task

	err := row.Scan(&task.Id, &task.Title, &task.CreatedAt, &task.DeletedAt, &task.SpentSeconds)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) AddTime(
	ctx context.Context,
	id int,
	seconds int,
) error {
	query := `
		UPDATE tasks
		SET spent_seconds = spent_seconds + $1
		WHERE id = $2 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query, seconds, id)
	return err
}

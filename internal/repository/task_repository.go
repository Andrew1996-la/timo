package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Andrew1996-la/timo/internal/models"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

const taskFields = `
	id,
	title,
	created_at,
	deleted_at,
	spent_seconds
`

func (r *TaskRepository) Create(ctx context.Context, title string) (*models.Task, error) {
	const query = `
		INSERT INTO tasks (title, created_at)
		VALUES (?, CURRENT_TIMESTAMP)
	`

	result, err := r.db.ExecContext(ctx, query, title)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get created task id: %w", err)
	}

	task, err := r.GetByID(ctx, int(id))
	if err != nil {
		return nil, fmt.Errorf("get created task: %w", err)
	}

	return task, nil
}

func (r *TaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM tasks
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`, taskFields)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)

	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id int) (*models.Task, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM tasks
		WHERE id = ? AND deleted_at IS NULL
	`, taskFields)

	task, err := scanTask(r.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}

		return nil, fmt.Errorf("get task id=%d: %w", id, err)
	}

	return &task, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int) error {
	const query = `
		UPDATE tasks
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete task id=%d: %w", id, err)
	}

	return checkRowsAffected(result, ErrTaskNotFound)
}

func (r *TaskRepository) AddTime(ctx context.Context, id int, seconds int) error {
	const query = `
		UPDATE tasks
		SET spent_seconds = spent_seconds + ?
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, seconds, id)
	if err != nil {
		return fmt.Errorf("add time to task id=%d: %w", id, err)
	}

	return checkRowsAffected(result, ErrTaskNotFound)
}

func checkRowsAffected(result sql.Result, notFoundErr error) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return notFoundErr
	}

	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanTask(s scanner) (models.Task, error) {
	var task models.Task

	err := s.Scan(
		&task.Id,
		&task.Title,
		&task.CreatedAt,
		&task.DeletedAt,
		&task.SpentSeconds,
	)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}
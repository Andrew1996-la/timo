package service

import (
	"context"
	"errors"

	"github.com/Andrew1996-la/timo/internal/models"
	"github.com/Andrew1996-la/timo/internal/repository"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (t *TaskService) Create(ctx context.Context, title string) (*models.Task, error) {
	if title == "" {
		return nil, errors.New("task title can not be empty")
	}

	return t.repo.Create(ctx, title)
}

func (t *TaskService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("task by id not found")
	}

	return t.repo.Delete(ctx, id)
}

func (t *TaskService) GetAll(ctx context.Context) ([]models.Task, error) {
	return t.repo.GetAll(ctx)
}

func (t *TaskService) GetById(ctx context.Context, id int) (*models.Task, error) {
	if id <= 0 {
		return nil, errors.New("task by id not found")
	}

	return t.repo.GetByID(ctx, id)
}

func (t *TaskService) AddTime(
	ctx context.Context,
	id int,
	seconds int,
) error {
	if seconds <= 0 {
		return nil
	}
	return t.repo.AddTime(ctx, id, seconds)
}

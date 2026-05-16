package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Andrew1996-la/timo/internal/models"
	"github.com/Andrew1996-la/timo/internal/repository"
)

var (
	ErrInvalidTaskID   = errors.New("invalid task id")
	ErrEmptyTaskTitle  = errors.New("task title cannot be empty")
	ErrInvalidDuration = errors.New("seconds must be positive")
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) Create(
	ctx context.Context,
	title string,
) (*models.Task, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return nil, ErrEmptyTaskTitle
	}

	return s.repo.Create(ctx, title)
}

func (s *TaskService) Delete(ctx context.Context, id int) error {
	if err := validateTaskID(id); err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}

func (s *TaskService) GetAll(ctx context.Context) ([]models.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *TaskService) GetByID(
	ctx context.Context,
	id int,
) (*models.Task, error) {
	if err := validateTaskID(id); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) AddTime(
	ctx context.Context,
	id int,
	seconds int,
) error {
	if err := validateTaskID(id); err != nil {
		return err
	}

	if seconds <= 0 {
		return ErrInvalidDuration
	}

	return s.repo.AddTime(ctx, id, seconds)
}

func validateTaskID(id int) error {
	if id <= 0 {
		return ErrInvalidTaskID
	}

	return nil
}
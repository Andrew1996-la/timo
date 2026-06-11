package handler

import (
	"context"

	"github.com/Andrew1996-la/timo/internal/models"
	"github.com/stretchr/testify/mock"
)

type mockTaskService struct {
	mock.Mock
}

func (m *mockTaskService) Create(
	ctx context.Context,
	title string,
) (*models.Task, error) {
	args := m.Called(ctx, title)

	task, _ := args.Get(0).(*models.Task)

	return task, args.Error(1)
}

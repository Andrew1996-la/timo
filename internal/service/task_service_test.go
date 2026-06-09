package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Andrew1996-la/timo/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockTaskRepository struct {
	mock.Mock
}

func (m *mockTaskRepository) Create(ctx context.Context, title string) (*models.Task, error) {
	args := m.Called(ctx, title)

	task, _ := args.Get(0).(*models.Task)

	return task, args.Error(1)
}

func (m *mockTaskRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func (m *mockTaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	args := m.Called(ctx)

	tasks, _ := args.Get(0).([]models.Task)

	return tasks, args.Error(1)
}

func (m *mockTaskRepository) GetByID(ctx context.Context, id int) (*models.Task, error) {
	args := m.Called(ctx, id)

	task, _ := args.Get(0).(*models.Task)

	return task, args.Error(1)
}

func (m *mockTaskRepository) AddTime(ctx context.Context, id int, seconds int) error {
	args := m.Called(ctx, id, seconds)

	return args.Error(0)
}

func TestTaskService_Create_EmptyTitle(t *testing.T) {
	repo := new(mockTaskRepository)
	taskService := NewTaskService(repo)

	task, err := taskService.Create(context.Background(), "")

	require.Error(t, err)
	assert.Nil(t, task)
	assert.True(t, errors.Is(err, ErrEmptyTaskTitle))

	repo.AssertNotCalled(t, "Create")
}

func TestTaskService_Create_OnlySpaces(t *testing.T) {
	repo := new(mockTaskRepository)
	taskService := NewTaskService(repo)

	task, err := taskService.Create(context.Background(), "     ")

	require.Error(t, err)
	assert.Nil(t, task)
	assert.True(t, errors.Is(err, ErrEmptyTaskTitle))

	repo.AssertNotCalled(t, "Create")
}

func TestTaskService_Create_Success(t *testing.T) {
	ctx := context.Background()

	repo := new(mockTaskRepository)
	taskService := NewTaskService(repo)

	expectedTask := &models.Task{
		Id:    1,
		Title: "Learn Go tests",
	}

	repo.On("Create", ctx, "Learn Go tests").Return(expectedTask, nil).Once()

	task, err := taskService.Create(ctx, "  Learn Go tests  ")

	require.NoError(t, err)
	assert.Equal(t, expectedTask, task)

	repo.AssertExpectations(t)
}

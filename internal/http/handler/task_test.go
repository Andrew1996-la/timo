package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Andrew1996-la/timo/internal/models"
	"github.com/Andrew1996-la/timo/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockTaskService struct {
	mock.Mock
}

func (m *mockTaskService) GetAll(ctx context.Context) ([]models.Task, error) {
	args := m.Called(ctx)

	tasks, _ := args.Get(0).([]models.Task)

	return tasks, args.Error(1)
}

func (m *mockTaskService) Create(ctx context.Context, title string) (*models.Task, error) {
	args := m.Called(ctx, title)

	task, _ := args.Get(0).(*models.Task)

	return task, args.Error(1)
}

func (m *mockTaskService) GetByID(ctx context.Context, id int) (*models.Task, error) {
	args := m.Called(ctx, id)

	task, _ := args.Get(0).(*models.Task)

	return task, args.Error(1)
}

func (m *mockTaskService) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func (m *mockTaskService) AddTime(ctx context.Context, id int, seconds int) error {
	args := m.Called(ctx, id, seconds)

	return args.Error(0)
}

func TestTaskHandler_Create_InvalidBody(t *testing.T) {
	taskService := new(mockTaskService)
	taskHandler := NewTaskHandler(taskService)

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader("{invalid json"),
	)
	w := httptest.NewRecorder()

	taskHandler.Tasks(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid body")

	taskService.AssertNotCalled(t, "Create")
}

func TestTaskHandler_Create_EmptyTitle(t *testing.T) {
	taskService := new(mockTaskService)
	taskHandler := NewTaskHandler(taskService)

	taskService.
		On("Create", mock.Anything, "").
		Return((*models.Task)(nil), service.ErrEmptyTaskTitle).
		Once()

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(`{"title": ""}`),
	)
	w := httptest.NewRecorder()

	taskHandler.Tasks(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), service.ErrEmptyTaskTitle.Error())

	taskService.AssertExpectations(t)
}

func TestTaskHandler_Create_Success(t *testing.T) {
	taskService := new(mockTaskService)
	taskHandler := NewTaskHandler(taskService)

	expectedTask := &models.Task{
		Id:    1,
		Title: "Learn handlers",
	}

	taskService.
		On("Create", mock.Anything, "Learn handlers").
		Return(expectedTask, nil).
		Once()

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(`{"title": "Learn handlers"}`),
	)
	w := httptest.NewRecorder()

	taskHandler.Tasks(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"title":"Learn handlers"`)

	taskService.AssertExpectations(t)
}

func TestTaskHandler_GetByID_InvalidID(t *testing.T) {
	taskService := new(mockTaskService)
	taskHandler := NewTaskHandler(taskService)

	req := httptest.NewRequest(http.MethodGet, "/tasks/abc", nil)
	w := httptest.NewRecorder()

	taskHandler.TaskByID(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid task id")

	taskService.AssertNotCalled(t, "GetByID")
}

func TestTaskHandler_Delete_Success(t *testing.T) {
	taskService := new(mockTaskService)
	taskHandler := NewTaskHandler(taskService)

	taskService.
		On("Delete", mock.Anything, 1).
		Return(nil).
		Once()

	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	w := httptest.NewRecorder()

	taskHandler.TaskByID(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)

	taskService.AssertExpectations(t)
}

func TestTaskHandler_AddTime_InvalidSeconds(t *testing.T) {
	taskService := new(mockTaskService)
	taskHandler := NewTaskHandler(taskService)

	taskService.
		On("AddTime", mock.Anything, 1, 0).
		Return(service.ErrInvalidDuration).
		Once()

	req := httptest.NewRequest(
		http.MethodPatch,
		"/tasks/1/time",
		strings.NewReader(`{"seconds": 0}`),
	)
	w := httptest.NewRecorder()

	taskHandler.AddTime(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), service.ErrInvalidDuration.Error())

	taskService.AssertExpectations(t)
}

func TestTaskHandler_ServiceInternalError(t *testing.T) {
	taskService := new(mockTaskService)
	taskHandler := NewTaskHandler(taskService)

	taskService.
		On("GetAll", mock.Anything).
		Return([]models.Task(nil), errors.New("db is down")).
		Once()

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	taskHandler.Tasks(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to get tasks")

	taskService.AssertExpectations(t)
}

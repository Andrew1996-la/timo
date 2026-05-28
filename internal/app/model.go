package app

import (
	"context"
	"time"

	"github.com/Andrew1996-la/timo/internal/models"
	"github.com/charmbracelet/bubbles/textinput"
)

type taskService interface {
	GetAll(ctx context.Context) ([]models.Task, error)
	Create(ctx context.Context, title string) (*models.Task, error)
	Delete(ctx context.Context, id int) error
	AddTime(ctx context.Context, id int, seconds int) error
}

type ViewMode int

const (
	ViewList ViewMode = iota
	ViewCreate
	ViewConfirmDelete
	ViewStats
)

type Model struct {
	Ctx     context.Context
	Service taskService

	Tasks    []models.Task
	Selected int
	Mode     ViewMode

	ConfirmDeleteTaskID int

	TimerRunning bool
	TimerTaskID  int
	TimerStarted time.Time

	Input textinput.Model
	Err   error
}

type tasksLoadedMsg struct {
	tasks []models.Task
	err   error
}

type taskCreatedMsg struct {
	err error
}

type taskDeletedMsg struct {
	err error
}

type timeAddedMsg struct {
	err error
}

type tickMsg time.Time
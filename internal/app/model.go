package app

import (
	"context"
	"time"

	"github.com/Andrew1996-la/timo/internal/models"
	"github.com/Andrew1996-la/timo/internal/service"
	"github.com/charmbracelet/bubbles/textinput"
)

type ViewMode int

const (
	ViewList ViewMode = iota
	ViewCreate
)

type Model struct {
	Ctx     context.Context
	Service *service.TaskService

	Tasks    []models.Task
	Selected int
	Mode     ViewMode

	TimerRunning bool
	TimerTaskID  int
	TimerStarted time.Time

	Input textinput.Model
	Err   error
}

// async messages
type tasksLoadedMsg struct {
	tasks []models.Task
	err   error
}
type taskCreatedMsg struct{ err error }
type taskDeletedMsg struct{ err error }
type tickMsg time.Time
type timeAddedMsg struct{ err error }

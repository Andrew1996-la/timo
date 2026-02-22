package app

import (
	"context"

	"github.com/Andrew1996-la/timo/internal/service"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func New(ctx context.Context, service *service.TaskService) tea.Model {
	ti := textinput.New()
	ti.Placeholder = "Task title"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 30

	return Model{
		Ctx:     ctx,
		Service: service,
		Mode:    ViewList,
		Input:   ti,
	}
}

func (m Model) Init() tea.Cmd {
	return loadTasks(m.Ctx, m.Service)
}

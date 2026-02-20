package app

import (
	"context"
	"fmt"

	"github.com/Andrew1996-la/timo/internal/models"
	"github.com/Andrew1996-la/timo/internal/service"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ctx     context.Context
	service *service.TaskService

	tasks    []models.Task
	selected int
	err      error
}

type tasksLoadedMsg struct {
	tasks []models.Task
	err   error
}

func New(ctx context.Context, service *service.TaskService) tea.Model {
	return Model{
		ctx:     ctx,
		service: service,
	}
}

func (m Model) Init() tea.Cmd {
	return loadTasks(m.ctx, m.service)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tasksLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.tasks = msg.tasks
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.selected > 0 {
				m.selected--
			}
		case "down":
			if m.selected < len(m.tasks)-1 {
				m.selected++
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n\npress q to quit"
	}

	if len(m.tasks) == 0 {
		return "No tasks yet.\n\npress q to quit"
	}

	result := "Tasks:\n\n"

	for i, t := range m.tasks {
		cursor := "  "
		if i == m.selected {
			cursor = "> "
		}
		result += fmt.Sprintf("%s [ID: %d] %s\n", cursor, t.Id, t.Title)
	}

	result += "\n↑↓ select   q quit"
	return result
}

func loadTasks(ctx context.Context, s *service.TaskService) tea.Cmd {
	return func() tea.Msg {
		tasks, err := s.GetAll(ctx)
		return tasksLoadedMsg{
			tasks: tasks,
			err:   err,
		}
	}
}

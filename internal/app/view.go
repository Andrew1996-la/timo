package app

import (
	"fmt"

	"github.com/Andrew1996-la/timo/internal/ui"
)

func (m Model) View() string {
	if m.Err != nil {
		return m.renderError()
	}

	switch m.Mode {
	case ViewCreate:
		return m.renderCreateView()
	case ViewConfirmDelete:
		return m.renderConfirmDeleteView()
	default:
		return m.renderListView()
	}
}

func (m Model) renderError() string {
	return fmt.Sprintf(
		"Error: %s\n\npress q to quit",
		m.Err.Error(),
	)
}

func (m Model) renderCreateView() string {
	return fmt.Sprintf(
		"Create task:\n\n%s\n\nenter save • esc cancel",
		m.Input.View(),
	)
}

func (m Model) renderConfirmDeleteView() string {
	return fmt.Sprintf(
		"Delete task:\n\n❗ %s\n\n[y] delete   [n] cancel",
		m.confirmDeleteTaskTitle(),
	)
}

func (m Model) renderListView() string {
	if len(m.Tasks) == 0 {
		return m.renderEmptyListView()
	}

	return ui.RenderTaskList(ui.TaskView{
		Tasks:        m.Tasks,
		Selected:     m.Selected,
		TimerRunning: m.TimerRunning,
		TimerTaskID:  m.TimerTaskID,
		TimerStarted: m.TimerStarted,
	})
}

func (m Model) renderEmptyListView() string {
	return "No tasks yet.\n\n↑↓ select   enter start/pause   n new   d delete   q quit"
}

func (m Model) confirmDeleteTaskTitle() string {
	for _, task := range m.Tasks {
		if task.Id == m.confirmDeleteTaskID {
			return task.Title
		}
	}

	return "unknown task"
}
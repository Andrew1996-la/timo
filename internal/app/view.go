package app

import (
	"fmt"

	"github.com/Andrew1996-la/timo/internal/ui"
)

const (
	quitHelp       = "press q to quit"
	createHelp     = "enter save • esc cancel"
	confirmHelp    = "[y] delete   [n] cancel"
	emptyListHelp  = "↑↓ select   enter start/pause   n new   d delete   q quit"
	unknownTask    = "unknown task"
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
	return fmt.Sprintf("Error: %s\n\n%s", m.Err.Error(), quitHelp)
}

func (m Model) renderCreateView() string {
	return fmt.Sprintf("Create task:\n\n%s\n\n%s", m.Input.View(), createHelp)
}

func (m Model) renderConfirmDeleteView() string {
	return fmt.Sprintf(
		"Delete task:\n\n❗ %s\n\n%s",
		m.confirmDeleteTaskTitle(),
		confirmHelp,
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
	return "No tasks yet.\n\n" + emptyListHelp
}

func (m Model) confirmDeleteTaskTitle() string {
	for _, task := range m.Tasks {
		if task.Id == m.confirmDeleteTaskID {
			return task.Title
		}
	}

	return unknownTask
}
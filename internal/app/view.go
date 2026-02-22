package app

import (
	"fmt"

	"github.com/Andrew1996-la/timo/internal/ui"
)

func (m Model) View() string {
	if m.Err != nil {
		return "Error: " + m.Err.Error() + "\n\npress q to quit"
	}

	if m.Mode == ViewCreate {
		return fmt.Sprintf("Create task:\n\n%s\n\nenter save • esc cancel", m.Input.View())
	}

	if len(m.Tasks) == 0 {
		return "No tasks yet.\n\npress q to quit"
	}

	// UI отдельной функции
	return ui.RenderTaskList(ui.TaskView{
		Tasks:        m.Tasks,
		Selected:     m.Selected,
		TimerRunning: m.TimerRunning,
		TimerTaskID:  m.TimerTaskID,
		TimerStarted: m.TimerStarted,
	})
}

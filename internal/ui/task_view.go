package ui

import (
	"fmt"
	"time"

	"github.com/Andrew1996-la/timo/internal/models"
)

// TaskView — данные для отображения задач
type TaskView struct {
	Tasks        []models.Task
	Selected     int
	TimerRunning bool
	TimerTaskID  int
	TimerStarted time.Time
}

// RenderTaskList — строит строку с задачами для терминала
func RenderTaskList(v TaskView) string {
	result := "Tasks:\n\n"

	for i, t := range v.Tasks {
		cursor := "  "
		if i == v.Selected {
			cursor = "> "
		}

		isRunning := v.TimerRunning && v.TimerTaskID == t.Id
		statusIcon := "⏸"
		if isRunning {
			statusIcon = "▶"
		}

		seconds := t.SpentSeconds
		if isRunning {
			seconds += int(time.Since(v.TimerStarted).Seconds())
		}

		line := fmt.Sprintf("%s %s %-20s %s", cursor, statusIcon, t.Title, formatSeconds(seconds))
		if isRunning {
			line += " ⚙️"
		}

		result += line + "\n"
	}

	result += "\n↑↓ select   enter start/pause   n new   d delete   q quit"
	return result
}

func formatSeconds(sec int) string {
	h := sec / 3600
	m := (sec % 3600) / 60
	s := sec % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

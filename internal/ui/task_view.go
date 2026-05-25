package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Andrew1996-la/timo/internal/models"
)

const taskTitleWidth = 20

type TaskView struct {
	Tasks        []models.Task
	Selected     int
	TimerRunning bool
	TimerTaskID  int
	TimerStarted time.Time
}

func RenderTaskList(v TaskView) string {
	var builder strings.Builder

	builder.WriteString("Tasks:\n\n")

	for i, task := range v.Tasks {
		builder.WriteString(renderTaskLine(v, i, task))
		builder.WriteString("\n")
	}

	builder.WriteString("\n↑↓ select   enter start/pause   n new   d delete   q quit")

	return builder.String()
}

func renderTaskLine(v TaskView, index int, task models.Task) string {
	isRunning := isTaskRunning(v, task)

	seconds := task.SpentSeconds
	if isRunning {
		seconds += int(time.Since(v.TimerStarted).Seconds())
	}

	line := fmt.Sprintf(
		"%s %s %-*s %s",
		cursor(index, v.Selected),
		statusIcon(isRunning),
		taskTitleWidth,
		task.Title,
		formatSeconds(seconds),
	)

	if isRunning {
		line += " ⚙️"
	}

	return line
}

func isTaskRunning(v TaskView, task models.Task) bool {
	return v.TimerRunning && v.TimerTaskID == task.Id
}

func cursor(index int, selected int) string {
	if index == selected {
		return ">"
	}

	return " "
}

func statusIcon(isRunning bool) string {
	if isRunning {
		return "▶"
	}

	return "⏸"
}

func formatSeconds(seconds int) string {
	hours := seconds / 3600
	minutes := seconds % 3600 / 60
	secs := seconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}
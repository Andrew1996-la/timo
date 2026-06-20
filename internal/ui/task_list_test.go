package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/Andrew1996-la/timo/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestRenderTaskList_ShowsTasks(t *testing.T) {
	view := TaskView{
		Tasks: []models.Task{
			{
				Id:           1,
				Title:        "Learn Go",
				SpentSeconds: 65,
			},
		},
		Selected: 0,
	}

	result := RenderTaskList(view)

	assert.Contains(t, result, "Tasks:")
	assert.Contains(t, result, "Learn Go")
	assert.Contains(t, result, "00:01:05")
}

func TestRenderTaskList_ShowsSelectedTask(t *testing.T) {
	view := TaskView{
		Tasks: []models.Task{
			{
				Id:    1,
				Title: "First task",
			},
			{
				Id:    2,
				Title: "Second task",
			},
		},
		Selected: 1,
	}

	result := RenderTaskList(view)

	lines := strings.Split(result, "\n")

	assert.Contains(t, lines[3], ">")
	assert.Contains(t, lines[3], "Second task")
}

func TestRenderTaskList_ShowsRunningTimer(t *testing.T) {
	started := time.Now().Add(-10 * time.Second)

	view := TaskView{
		Tasks: []models.Task{
			{
				Id:           1,
				Title:        "Running task",
				SpentSeconds: 60,
			},
		},
		Selected:     0,
		TimerRunning: true,
		TimerTaskID:  1,
		TimerStarted: started,
	}

	result := RenderTaskList(view)

	assert.Contains(t, result, "Running task")
	assert.Contains(t, result, "▶")
	assert.Contains(t, result, "⚙️")
}

func TestFormatSeconds(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int
		expected string
	}{
		{
			name:     "zero",
			seconds:  0,
			expected: "00:00:00",
		},
		{
			name:     "seconds",
			seconds:  7,
			expected: "00:00:07",
		},
		{
			name:     "minutes and seconds",
			seconds:  65,
			expected: "00:01:05",
		},
		{
			name:     "hours minutes seconds",
			seconds:  3661,
			expected: "01:01:01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, formatSeconds(tt.seconds))
		})
	}
}
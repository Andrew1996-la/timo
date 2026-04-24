package app

import (
	"context"
	"time"

	"github.com/Andrew1996-la/timo/internal/service"
	tea "github.com/charmbracelet/bubbletea"
)

func loadTasks(ctx context.Context, taskService *service.TaskService) tea.Cmd {
	return func() tea.Msg {
		tasks, err := taskService.GetAll(ctx)

		return tasksLoadedMsg{
			tasks: tasks,
			err:   err,
		}
	}
}

func createTask(ctx context.Context, taskService *service.TaskService, title string) tea.Cmd {
	return func() tea.Msg {
		_, err := taskService.Create(ctx, title)

		return taskCreatedMsg{
			err: err,
		}
	}
}

func deleteTask(ctx context.Context, taskService *service.TaskService, id int) tea.Cmd {
	return func() tea.Msg {
		err := taskService.Delete(ctx, id)

		return taskDeletedMsg{
			err: err,
		}
	}
}

func addTime(ctx context.Context, taskService *service.TaskService, id int, seconds int) tea.Cmd {
	return func() tea.Msg {
		err := taskService.AddTime(ctx, id, seconds)

		return timeAddedMsg{
			err: err,
		}
	}
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
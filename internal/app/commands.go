package app

import (
	"context"
	"time"

	"github.com/Andrew1996-la/timo/internal/service"
	tea "github.com/charmbracelet/bubbletea"
)

func loadTasks(ctx context.Context, s *service.TaskService) tea.Cmd {
	return func() tea.Msg {
		tasks, err := s.GetAll(ctx)
		return tasksLoadedMsg{tasks: tasks, err: err}
	}
}

func createTask(ctx context.Context, s *service.TaskService, title string) tea.Cmd {
	return func() tea.Msg {
		_, err := s.Create(ctx, title)
		return taskCreatedMsg{err: err}
	}
}

func deleteTask(ctx context.Context, s *service.TaskService, id int) tea.Cmd {
	return func() tea.Msg {
		err := s.Delete(ctx, id)
		return taskDeletedMsg{err: err}
	}
}

func addTime(ctx context.Context, s *service.TaskService, id int, seconds int) tea.Cmd {
	return func() tea.Msg {
		err := s.AddTime(ctx, id, seconds)
		return timeAddedMsg{err: err}
	}
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

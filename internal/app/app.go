package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Andrew1996-la/timo/internal/models"
	"github.com/Andrew1996-la/timo/internal/service"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewMode int

const (
	ViewList ViewMode = iota
	ViewCreate
)

type Model struct {
	ctx     context.Context
	service *service.TaskService

	tasks    []models.Task
	selected int
	mode     ViewMode

	timerRunning bool
	timerTaskID  int
	timerStarted time.Time

	input textinput.Model
	err   error
}

type tasksLoadedMsg struct {
	tasks []models.Task
	err   error
}

type taskCreatedMsg struct {
	err error
}

type taskDeletedMsg struct {
	err error
}

type tickMsg time.Time

type timeAddedMsg struct {
	err error
}

func New(ctx context.Context, service *service.TaskService) tea.Model {
	ti := textinput.New()
	ti.Placeholder = "Task title"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 30

	return Model{
		ctx:     ctx,
		service: service,
		mode:    ViewList,
		input:   ti,
	}
}

func (m Model) Init() tea.Cmd {
	return loadTasks(m.ctx, m.service)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// 1️⃣ СНАЧАЛА — async сообщения
	switch msg := msg.(type) {

	case tasksLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.tasks = msg.tasks
		return m, nil

	case taskCreatedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		m.mode = ViewList
		return m, loadTasks(m.ctx, m.service)

	case taskDeletedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		// корректируем selected, если удалили последний элемент
		if m.selected > 0 {
			m.selected--
		}

		return m, loadTasks(m.ctx, m.service)

	case tickMsg:
		if m.timerRunning {
			return m, tick()
		}
		return m, nil

	case timeAddedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		return m, loadTasks(m.ctx, m.service)
	}

	// 2️⃣ ПОТОМ — режим создания
	if m.mode == ViewCreate {
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch msg.Type {

			case tea.KeyEnter:
				title := m.input.Value()
				if title == "" {
					return m, nil
				}
				m.input.Blur()
				return m, createTask(m.ctx, m.service, title)

			case tea.KeyEsc:
				m.mode = ViewList
				m.input.Blur()
				return m, nil
			}
		}

		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	// 3️⃣ Обычный список
	switch msg := msg.(type) {

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

		case "n":
			m.mode = ViewCreate
			m.input.SetValue("")
			m.input.Focus()
			return m, nil

		case "enter":
			if len(m.tasks) == 0 {
				return m, nil
			}

			task := m.tasks[m.selected]

			// ⏹ если таймер уже бежит — останавливаем
			if m.timerRunning {
				elapsed := int(time.Since(m.timerStarted).Seconds())

				stopID := m.timerTaskID

				m.timerRunning = false
				m.timerTaskID = 0

				// сохраняем время
				cmd := addTime(m.ctx, m.service, stopID, elapsed)

				// если нажали ENTER на ТУ ЖЕ задаче — просто стоп
				if stopID == task.Id {
					return m, cmd
				}

				// иначе — старт новой
				m.timerRunning = true
				m.timerTaskID = task.Id
				m.timerStarted = time.Now()

				return m, tea.Batch(cmd, tick())
			}

			// ▶️ старт, если таймер не бежит
			m.timerRunning = true
			m.timerTaskID = task.Id
			m.timerStarted = time.Now()

			return m, tick()

		case "d":
			if len(m.tasks) == 0 {
				return m, nil
			}

			task := m.tasks[m.selected]
			return m, deleteTask(m.ctx, m.service, task.Id)

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

	if m.mode == ViewCreate {
		return fmt.Sprintf(
			"Create task:\n\n%s\n\nenter save • esc cancel",
			m.input.View(),
		)
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

		seconds := t.SpentSeconds

		// если таймер запущен на этой задаче — считаем live-время
		if m.timerRunning && m.timerTaskID == t.Id {
			seconds += int(time.Since(m.timerStarted).Seconds())
		}

		timeStr := formatSeconds(seconds)

		result += fmt.Sprintf(
			"%s %s   %s\n",
			cursor,
			t.Title,
			timeStr,
		)
	}

	result += "\n↑↓ select   n new   d delete   q quit"
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

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func addTime(ctx context.Context, s *service.TaskService, id int, seconds int) tea.Cmd {
	return func() tea.Msg {
		err := s.AddTime(ctx, id, seconds)
		return timeAddedMsg{err: err}
	}
}

func formatSeconds(sec int) string {
	h := sec / 3600
	m := (sec % 3600) / 60
	s := sec % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

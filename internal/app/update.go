package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Асинхронные сообщения
	switch msg := msg.(type) {
	case tasksLoadedMsg:
		return m.updateTasksLoaded(msg)
	case taskCreatedMsg:
		return m.updateTaskCreated(msg)
	case taskDeletedMsg:
		return m.updateTaskDeleted(msg)
	case timeAddedMsg:
		return m.updateTimeAdded(msg)
	case tickMsg:
		return m.updateTick()
	}

	// Режим создания новой задачи
	if m.Mode == ViewCreate {
		return m.updateCreateMode(msg)
	}

	// Режим списка задач (основной)
	return m.updateListMode(msg)
}

func (m Model) updateTasksLoaded(msg tasksLoadedMsg) (Model, tea.Cmd) {
	if msg.err != nil {
		m.Err = msg.err
		return m, nil
	}
	m.Tasks = msg.tasks
	return m, nil
}

func (m Model) updateTaskCreated(msg taskCreatedMsg) (Model, tea.Cmd) {
	if msg.err != nil {
		m.Err = msg.err
		return m, nil
	}
	m.Mode = ViewList
	return m, loadTasks(m.Ctx, m.Service)
}

func (m Model) updateTaskDeleted(msg taskDeletedMsg) (Model, tea.Cmd) {
	if msg.err != nil {
		m.Err = msg.err
		return m, nil
	}

	// корректируем selected, если удалили последний элемент
	if m.Selected > 0 {
		m.Selected--
	}

	return m, loadTasks(m.Ctx, m.Service)
}

func (m Model) updateTimeAdded(msg timeAddedMsg) (Model, tea.Cmd) {
	if msg.err != nil {
		m.Err = msg.err
		return m, nil
	}
	return m, loadTasks(m.Ctx, m.Service)
}

func (m Model) updateTick() (Model, tea.Cmd) {
	if m.TimerRunning {
		return m, tick()
	}
	return m, nil
}

func (m Model) updateCreateMode(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			title := m.Input.Value()
			if title == "" {
				return m, nil
			}
			m.Input.Blur()
			return m, createTask(m.Ctx, m.Service, title)
		case tea.KeyEsc:
			m.Mode = ViewList
			m.Input.Blur()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m Model) updateListMode(msg tea.Msg) (Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.String() {
	case "up":
		if m.Selected > 0 {
			m.Selected--
		}
	case "down":
		if m.Selected < len(m.Tasks)-1 {
			m.Selected++
		}
	case "n":
		m.Mode = ViewCreate
		m.Input.SetValue("")
		m.Input.Focus()
	case "enter":
		return m.toggleTimer()
	case "d":
		if len(m.Tasks) == 0 {
			return m, nil
		}
		task := m.Tasks[m.Selected]
		return m, deleteTask(m.Ctx, m.Service, task.Id)
	case "q", "ctrl+c":
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) toggleTimer() (Model, tea.Cmd) {
	if len(m.Tasks) == 0 {
		return m, nil
	}

	task := m.Tasks[m.Selected]

	// Если таймер бежит - ставим на паузу
	if m.TimerRunning {
		elapsed := int(time.Since(m.TimerStarted).Seconds())
		stopID := m.TimerTaskID

		m.TimerRunning = false
		m.TimerTaskID = 0

		cmd := addTime(m.Ctx, m.Service, stopID, elapsed)

		// Если нажали Enter на ту же задачу - просто пауза
		if stopID == task.Id {
			return m, cmd
		}

		// Иначе - старт новой задачи
		m.TimerRunning = true
		m.TimerTaskID = task.Id
		m.TimerStarted = time.Now()

		return m, tea.Batch(cmd, tick())
	}

	// Таймер не бежит - стартуем
	m.TimerRunning = true
	m.TimerTaskID = task.Id
	m.TimerStarted = time.Now()

	return m, tick()
}

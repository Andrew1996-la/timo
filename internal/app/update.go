package app

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if updated, cmd, ok := m.handleAsyncMsg(msg); ok {
		return updated, cmd
	}

	switch m.Mode {
	case ViewCreate:
		return m.updateCreateMode(msg)
	case ViewConfirmDelete:
		return m.updateConfirmDeleteMode(msg)
	default:
		return m.updateListMode(msg)
	}
}

func (m Model) handleAsyncMsg(msg tea.Msg) (Model, tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tasksLoadedMsg:
		return m.updateTasksLoaded(msg), nil, true
	case taskCreatedMsg:
		updated, cmd := m.updateTaskCreated(msg)
		return updated, cmd, true
	case taskDeletedMsg:
		updated, cmd := m.updateTaskDeleted(msg)
		return updated, cmd, true
	case timeAddedMsg:
		updated, cmd := m.updateTimeAdded(msg)
		return updated, cmd, true
	case tickMsg:
		updated, cmd := m.updateTick()
		return updated, cmd, true
	default:
		return m, nil, false
	}
}

func (m Model) updateTasksLoaded(msg tasksLoadedMsg) Model {
	if msg.err != nil {
		return m.withError(msg.err)
	}

	m.Tasks = msg.tasks
	m = m.normalizeSelected()

	return m
}

func (m Model) updateTaskCreated(msg taskCreatedMsg) (Model, tea.Cmd) {
	if msg.err != nil {
		return m.withError(msg.err), nil
	}

	m.Mode = ViewList
	m.Input.SetValue("")
	m.Input.Blur()

	return m, loadTasks(m.Ctx, m.Service)
}

func (m Model) updateTaskDeleted(msg taskDeletedMsg) (Model, tea.Cmd) {
	if msg.err != nil {
		return m.withError(msg.err), nil
	}

	m.Mode = ViewList
	m.confirmDeleteTaskID = 0

	return m, loadTasks(m.Ctx, m.Service)
}

func (m Model) updateTimeAdded(msg timeAddedMsg) (Model, tea.Cmd) {
	if msg.err != nil {
		return m.withError(msg.err), nil
	}

	return m, loadTasks(m.Ctx, m.Service)
}

func (m Model) updateTick() (Model, tea.Cmd) {
	if !m.TimerRunning {
		return m, nil
	}

	return m, tick()
}

func (m Model) updateCreateMode(msg tea.Msg) (Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		var cmd tea.Cmd
		m.Input, cmd = m.Input.Update(msg)
		return m, cmd
	}

	switch key.Type {
	case tea.KeyEnter:
		return m.submitCreateTask()
	case tea.KeyEsc:
		return m.cancelCreateTask()
	default:
		var cmd tea.Cmd
		m.Input, cmd = m.Input.Update(msg)
		return m, cmd
	}
}

func (m Model) submitCreateTask() (Model, tea.Cmd) {
	title := strings.TrimSpace(m.Input.Value())
	if title == "" {
		return m, nil
	}

	m.Input.Blur()

	return m, createTask(m.Ctx, m.Service, title)
}

func (m Model) cancelCreateTask() (Model, tea.Cmd) {
	m.Mode = ViewList
	m.Input.SetValue("")
	m.Input.Blur()

	return m, nil
}

func (m Model) updateConfirmDeleteMode(msg tea.Msg) (Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.String() {
	case "y":
		return m.confirmDelete()
	case "n", "esc":
		return m.cancelDelete()
	default:
		return m, nil
	}
}

func (m Model) confirmDelete() (Model, tea.Cmd) {
	id := m.confirmDeleteTaskID

	m.confirmDeleteTaskID = 0
	m.Mode = ViewList

	return m, deleteTask(m.Ctx, m.Service, id)
}

func (m Model) cancelDelete() (Model, tea.Cmd) {
	m.confirmDeleteTaskID = 0
	m.Mode = ViewList

	return m, nil
}

func (m Model) updateListMode(msg tea.Msg) (Model, tea.Cmd) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch key.String() {
	case "up":
		return m.moveSelectionUp(), nil
	case "down":
		return m.moveSelectionDown(), nil
	case "n":
		return m.openCreateMode()
	case "enter":
		return m.toggleTimer()
	case "d":
		return m.openDeleteConfirm()
	case "q", "ctrl+c":
		return m, tea.Quit
	default:
		return m, nil
	}
}

func (m Model) moveSelectionUp() Model {
	if m.Selected > 0 {
		m.Selected--
	}

	return m
}

func (m Model) moveSelectionDown() Model {
	if m.Selected < len(m.Tasks)-1 {
		m.Selected++
	}

	return m
}

func (m Model) openCreateMode() (Model, tea.Cmd) {
	m.Mode = ViewCreate
	m.Input.SetValue("")
	m.Input.Focus()

	return m, nil
}

func (m Model) openDeleteConfirm() (Model, tea.Cmd) {
	if len(m.Tasks) == 0 {
		return m, nil
	}

	task := m.Tasks[m.Selected]

	m.Mode = ViewConfirmDelete
	m.confirmDeleteTaskID = task.Id

	return m, nil
}

func (m Model) toggleTimer() (Model, tea.Cmd) {
	if len(m.Tasks) == 0 {
		return m, nil
	}

	selectedTask := m.Tasks[m.Selected]

	if !m.TimerRunning {
		return m.startTimer(selectedTask.Id)
	}

	return m.stopOrSwitchTimer(selectedTask.Id)
}

func (m Model) startTimer(taskID int) (Model, tea.Cmd) {
	m.TimerRunning = true
	m.TimerTaskID = taskID
	m.TimerStarted = time.Now()

	return m, tick()
}

func (m Model) stopOrSwitchTimer(selectedTaskID int) (Model, tea.Cmd) {
	elapsed := int(time.Since(m.TimerStarted).Seconds())
	stoppedTaskID := m.TimerTaskID

	m.TimerRunning = false
	m.TimerTaskID = 0

	saveTimeCmd := addTime(m.Ctx, m.Service, stoppedTaskID, elapsed)

	if stoppedTaskID == selectedTaskID {
		return m, saveTimeCmd
	}

	m.TimerRunning = true
	m.TimerTaskID = selectedTaskID
	m.TimerStarted = time.Now()

	return m, tea.Batch(saveTimeCmd, tick())
}

func (m Model) normalizeSelected() Model {
	if len(m.Tasks) == 0 {
		m.Selected = 0
		return m
	}

	if m.Selected >= len(m.Tasks) {
		m.Selected = len(m.Tasks) - 1
	}

	if m.Selected < 0 {
		m.Selected = 0
	}

	return m
}

func (m Model) withError(err error) Model {
	m.Err = err
	return m
}
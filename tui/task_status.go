package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evangodon/todomd/task"
)

type statusAction int

const (
	nextStatusType statusAction = iota
	prevStatusType
)

type taskActionFunc func(m model) model

type TaskActionMsg struct {
	action         statusAction
	task           task.Task
	runAfterUpdate taskActionFunc
}

func newTaskActionType(action statusAction, task task.Task, runAfterUpdate taskActionFunc) tea.Cmd {
	return func() tea.Msg {
		return TaskActionMsg{
			action:         action,
			task:           task,
			runAfterUpdate: runAfterUpdate,
		}
	}
}

func (m model) handleTaskAction(msg TaskActionMsg) (model, tea.Cmd) {
	// check if todo is nil
	switch msg.action {
	case nextStatusType:
		for _, t := range m.todosList.Tasks() {
			if t.Body() == msg.task.Body() {
				t.SetStatus(t.Status().Next())
			}
		}
	case prevStatusType:
		for _, t := range m.todosList.Tasks() {
			if t.Body() == msg.task.Body() {
				t.SetStatus(t.Status().Prev())
			}
		}
	}

	m.groups = updateGroups(m.todosList)

	if msg.runAfterUpdate != nil {
		m = msg.runAfterUpdate(m)
	}

	return m, nil
}

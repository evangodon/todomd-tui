package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evangodon/todomd/task"
)

type actionType int

const (
	nextStatusType actionType = iota
	prevStatusType
)

// TODO: think of better names
type taskActionFunc func(m model) model

type TaskAction struct {
	name           actionType
	task           task.Task
	runAfterUpdate taskActionFunc
}

type TaskActionMsg TaskAction

func newTaskActionType(name actionType, task task.Task, runAfterUpdate taskActionFunc) tea.Cmd {

	return func() tea.Msg {
		return TaskActionMsg{
			name:           name,
			task:           task,
			runAfterUpdate: runAfterUpdate,
		}
	}
}

func (m model) handleTaskAction(msg TaskActionMsg) (model, tea.Cmd) {
	// check if todo is nil
	switch msg.name {
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

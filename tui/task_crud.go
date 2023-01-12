package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evangodon/todomd/task"
)

type taskAction int

const (
	addTask taskAction = iota
	removeTask
)

type TaskMsg struct {
	action taskAction
	task   *task.Task
}

func NewTaskMsg(action taskAction, task *task.Task) tea.Cmd {
	return func() tea.Msg {
		return TaskMsg{
			action: action,
			task:   task,
		}
	}
}

func (m model) handleNewTaskMsg(msg TaskMsg) (model, tea.Cmd) {

	switch msg.action {
	case addTask:
		m.todosList.AddTask(msg.task)
		if msg.task.Parent() != nil {
			msg.task.Parent().AddSubTask(msg.task)
		}
	case removeTask:
		m.todosList.RemoveTask(msg.task)
		if msg.task.Parent() != nil {
			msg.task.Parent().RemoveSubTask(msg.task)
		}
	}

	m.groups = updateGroups(m.todosList)
	return m, nil
}

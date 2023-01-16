package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evangodon/todomd/task"
)

type keyMap struct {
	movement       key.Binding
	left           key.Binding
	right          key.Binding
	down           key.Binding
	up             key.Binding
	moveNextStatus key.Binding
	prevNextStatus key.Binding
	startStop      key.Binding
	newTask        key.Binding
	deleteTask     key.Binding
	newSubTask     key.Binding
	help           key.Binding
	quit           key.Binding
}

var keys = keyMap{
	movement: key.NewBinding(
		key.WithKeys("←↑↓→"),
		key.WithHelp("hjkl", "navigate"),
	),
	left: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("h", "left"),
	),
	right: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("l", "right"),
	),
	up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k", "up"),
	),
	down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j", "down"),
	),
	startStop: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "start/stop"),
	),
	prevNextStatus: key.NewBinding(
		key.WithKeys("H", "shift+left"),
		key.WithHelp("H", "move todo left"),
	),
	moveNextStatus: key.NewBinding(
		key.WithKeys("L", "shift+right"),
		key.WithHelp("L", "move todo right"),
	),
	newTask: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add task"),
	),
	newSubTask: key.NewBinding(
		key.WithKeys("ctrl+a"),
		key.WithHelp("ctrl+a", "add subtask"),
	),
	deleteTask: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "delete task"),
	),
	help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.help, k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	k.help.SetHelp("?", "Close help")
	return [][]key.Binding{
		{k.movement, k.newTask},
		{k.moveNextStatus, k.prevNextStatus},
		{k.help, k.quit},
	}
}

// TODO: add method for getting selected task

func (m model) handleKey(msg tea.KeyMsg) (model, tea.Cmd) {
	activeGroup := m.activeGroup()

	switch {
	// LEFT
	case key.Matches(msg, keys.left):
		next := m.position.GoLeft()
		return m, newPositionMsg(next)

	// RIGHT
	case key.Matches(msg, keys.right):
		m.position = m.position.GoRight()
		return m, newPositionMsg(m.position)

	// UP
	case key.Matches(msg, keys.up):
		m.position = m.position.GoUp()
		return m, newPositionMsg(m.position)

	// DOWN
	case key.Matches(msg, keys.down):
		m.position = m.position.GoDown()
		return m, newPositionMsg(m.position)

	// START / STOP
	case key.Matches(msg, keys.startStop):
		todo := activeGroup.Tasks()[m.position.Y]

		var moveToExistingTask = func(m model) model {
			m.position.Y = task.Clamp(0, m.position.Y, len(m.activeGroup().Tasks())-1)
			return m
		}

		if todo.Status() == task.UncompletedStatus {
			return m, newTaskActionType(nextStatusType, todo, moveToExistingTask)
		}

		if todo.Status() == task.InProgressStatus {
			return m, newTaskActionType(prevStatusType, todo, moveToExistingTask)
		}

		return m, nil

	// PREVIOUS STATUS
	case key.Matches(msg, keys.prevNextStatus):
		todo := activeGroup.Tasks()[m.position.Y]

		return m, newTaskActionType(prevStatusType, todo, func(m model) model {
			nextPos := m.position.GoLeft()
			nextGroup := m.GetGroup(nextPos.X)
			for i, t := range nextGroup.Tasks() {
				if t.Body() == todo.Body() {
					nextPos.Y = i
					m.position = nextPos
					return m
				}
			}
			return m
		})

	// NEXT STATUS
	case key.Matches(msg, keys.moveNextStatus):
		todo := activeGroup.Tasks()[m.position.Y]

		return m, newTaskActionType(nextStatusType, todo, func(m model) model {
			nextPos := m.position.GoRight()
			nextGroup := m.GetGroup(nextPos.X)
			for i, t := range nextGroup.Tasks() {
				if t.Body() == todo.Body() {
					nextPos.Y = i
					m.position = nextPos
					return m
				}
			}
			return m
		})

		// ADD TASK
	case key.Matches(msg, keys.newTask):
		m.textinput.enabled = true
		m.textinput.input.Width = task.Clamp(10, activeGroup.Width()-4, 50)
		return m, tea.Batch(textinput.Blink, m.textinput.input.Focus())

		// REMOVE TASK
	case key.Matches(msg, keys.deleteTask):
		selectedTask := activeGroup.Tasks()[m.position.Y]
		var taskToDelete *task.Task
		for _, t := range m.todosList.Tasks() {
			if t.Body() == selectedTask.Body() {
				taskToDelete = t
			}
		}
		return m, NewTaskMsg(removeTask, taskToDelete)

		// ADD SUBTASK
	case key.Matches(msg, keys.newSubTask):
		selectedTask := activeGroup.Tasks()[m.position.Y]
		// Disable if task is already subtask
		if selectedTask.Parent() != nil {
			return m, nil
		}
		m.textinput.enabled = true
		for _, t := range m.todosList.Tasks() {
			if t.Body() == selectedTask.Body() {
				m.textinput.parent = t
			}
		}
		m.textinput.input.Width = task.Clamp(10, activeGroup.Width()-4, 50)
		return m, tea.Batch(textinput.Blink, m.textinput.input.Focus())

		// HELP
	case key.Matches(msg, keys.help):
		m.help.ShowAll = !m.help.ShowAll
		return m, nil
		// QUIT
	case key.Matches(msg, keys.quit):
		m.todosList.WriteToFile()
		return m, tea.Quit
	default:
		return m, nil
	}
}

func newKeyMsg(key string) tea.Cmd {
	return func() tea.Msg {
		return tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune(key),
		}
	}
}

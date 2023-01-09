package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evangodon/todomd/task"
)

type keyMap struct {
	movement   key.Binding
	left       key.Binding
	right      key.Binding
	down       key.Binding
	up         key.Binding
	nextStatus key.Binding
	prevStatus key.Binding
	add        key.Binding
	help       key.Binding
	quit       key.Binding
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
	prevStatus: key.NewBinding(
		key.WithKeys("H", "shift+left"),
		key.WithHelp("H", "move todo left"),
	),
	nextStatus: key.NewBinding(
		key.WithKeys("L", "shift+right"),
		key.WithHelp("L", "move todo right"),
	),
	add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add todo"),
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
		{k.movement, k.add},
		{k.nextStatus, k.prevStatus},
		{k.help, k.quit},
	}
}

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

	// PREVIOUS STATUS
	case key.Matches(msg, keys.prevStatus):
		todo := activeGroup.Tasks()[m.position.Y]
		for _, t := range m.todosList.Tasks() {
			if t.Body() == todo.Body() {
				t.SetStatus(t.Status().Prev())
			}
		}

		m.groups = updateGroups(m.todosList)
		nextPos := m.position.GoLeft()
		nextGroup := m.GetGroup(nextPos.X)

		for i, t := range nextGroup.Tasks() {
			if t.Body() == todo.Body() {
				nextPos.Y = i

				return m, newPositionMsg(nextPos)
			}
		}
		return m, newPositionMsg(Position{0, 0})

	// NEXT STATUS
	case key.Matches(msg, keys.nextStatus):
		todo := activeGroup.Tasks()[m.position.Y]
		for _, t := range m.todosList.Tasks() {
			if t.Body() == todo.Body() {
				t.SetStatus(t.Status().Next())
			}
		}

		m.groups = updateGroups(m.todosList)
		nextPos := m.position.GoRight()
		nextGroup := m.GetGroup(nextPos.X)

		for i, t := range nextGroup.Tasks() {
			if t.Body() == todo.Body() {
				nextPos.Y = i
				return m, newPositionMsg(nextPos)
			}
		}
		return m, newPositionMsg(Position{0, 0})

		// ADD TODO
	case key.Matches(msg, keys.add):
		m.textinput.enabled = true
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

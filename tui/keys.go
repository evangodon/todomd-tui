package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evangodon/todo/internal"
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
		key.WithHelp("?", "toggle help"),
	),
	quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.movement, k.add, k.nextStatus, k.prevStatus, k.help, k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.movement, k.add},
		{k.nextStatus, k.prevStatus},
		{k.help, k.quit},
	}
}

func (m model) handleKey(msg tea.KeyMsg) (model, tea.Cmd) {
	activeGroup := m.activeGroup()
	maxY := len(activeGroup.Items()) - 1

	switch {
	// LEFT
	case key.Matches(msg, keys.left):
		m.position.X = internal.Clamp(0, m.position.X-1, 2)
		activeGroup = m.activeGroup()
		items := activeGroup.Items()
		maxY := len(items) - 1
		if len(items) == 0 && m.position.X > 0 {
			return m, newKeyMsg("h")
		}
		m.position.Y = internal.Clamp(0, m.position.Y, maxY)
		return m, nil
	// RIGHT
	case key.Matches(msg, keys.right):
		m.position.X = internal.Clamp(0, m.position.X+1, 2)
		activeGroup = m.activeGroup()
		maxY := len(activeGroup.Items()) - 1
		if len(activeGroup.Items()) == 0 && m.position.X < 2 {
			return m, newKeyMsg("l")
		}
		m.position.Y = internal.Clamp(0, m.position.Y, maxY)
		return m, nil
	// UP
	case key.Matches(msg, keys.up):
		m.position.Y = internal.Clamp(0, m.position.Y-1, maxY)
		return m, nil
	// DOWN
	case key.Matches(msg, keys.down):
		m.position.Y = internal.Clamp(0, m.position.Y+1, maxY)
		return m, nil
	// PREVIOUS STATUS
	case key.Matches(msg, keys.prevStatus):
		todo := activeGroup.Items()[m.position.Y]
		for _, t := range m.todosList.Items() {
			if t.Body() == todo.Body() {
				t.Status = t.Status.Prev()
			}
		}

		m.groups = updateGroups(m.todosList)
		return m, func() tea.Msg {
			nextX := internal.Clamp(0, m.position.X-1, 2)
			m.position.X = nextX
			currGroup := m.activeGroup()

			for i, t := range currGroup.Items() {
				if t.Body() == todo.Body() {
					nextY := i

					return newPositionMsg(nextX, nextY)
				}
			}
			return newPositionMsg(0, 0)
		}
	// NEXT STATUS
	case key.Matches(msg, keys.nextStatus):
		todo := activeGroup.Items()[m.position.Y]
		for _, t := range m.todosList.Items() {
			if t.Body() == todo.Body() {
				t.Status = t.Status.Next()
			}
		}

		m.groups = updateGroups(m.todosList)
		return m, func() tea.Msg {
			nextX := internal.Clamp(0, m.position.X+1, 2)
			m.position.X = nextX
			currGroup := m.activeGroup()
			for i, t := range currGroup.Items() {
				if t.Body() == todo.Body() {
					nextY := i

					return newPositionMsg(nextX, nextY)
				}
			}
			return newPositionMsg(0, 0)
		}
		// ADD TODO
	case key.Matches(msg, keys.add):
		m.textinput.enabled = true
		m.textinput.input.Width = internal.Clamp(10, activeGroup.Width(), 50)
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

func newPositionMsg(x, y int) tea.Msg {
	return Position{
		X: x,
		Y: y,
	}
}

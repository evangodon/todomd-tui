package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
	movement   key.Binding
	left       key.Binding
	right      key.Binding
	down       key.Binding
	up         key.Binding
	nextStatus key.Binding
	prevStatus key.Binding
	Help       key.Binding
	Quit       key.Binding
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
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.movement, k.nextStatus, k.prevStatus, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.movement},
		{k.nextStatus, k.prevStatus},
		{k.Help, k.Quit},
	}
}

func (m model) handleKey(msg tea.KeyMsg) (model, tea.Cmd) {
	activeGroup := m.activeGroup()
	maxY := len(activeGroup.items) - 1
	switch {
	// LEFT
	case key.Matches(msg, keys.left):
		m.position.x = clamp(0, m.position.x-1, 2)
		activeGroup = m.activeGroup()
		maxY := len(activeGroup.items) - 1
		if len(activeGroup.items) == 0 && m.position.x > 0 {
			return m, newKeyMsg("h")
		}
		m.position.y = clamp(0, m.position.y, maxY)
		return m, nil
	// RIGHT
	case key.Matches(msg, keys.right):
		m.position.x = clamp(0, m.position.x+1, 2)
		activeGroup = m.activeGroup()
		maxY := len(activeGroup.items) - 1
		if len(activeGroup.items) == 0 && m.position.x < 2 {
			return m, newKeyMsg("l")
		}
		m.position.y = clamp(0, m.position.y, maxY)
		return m, nil
	// UP
	case key.Matches(msg, keys.up):
		m.position.y = clamp(0, m.position.y-1, maxY)
		return m, nil
	// DOWN
	case key.Matches(msg, keys.down):
		m.position.y = clamp(0, m.position.y+1, maxY)
		return m, nil
	// PREVIOUS STATUS
	case key.Matches(msg, keys.prevStatus):
		todo := activeGroup.items[m.position.y]
		for _, t := range m.todosList.items {
			if t.body == todo.body {
				t.status = t.status.prev()
			}
		}

		m.groups = updateGroups(m.todosList)
		return m, func() tea.Msg {
			nextX := clamp(0, m.position.x-1, maxY)
			m.position.x = nextX
			currGroup := m.activeGroup()

			for i, t := range currGroup.items {
				if t.body == todo.body {
					nextY := i

					return newPositionMsg(nextX, nextY)
				}
			}
			return newPositionMsg(0, 0)
		}
	// NEXT STATUS
	case key.Matches(msg, keys.nextStatus):
		todo := activeGroup.items[m.position.y]
		for _, t := range m.todosList.items {
			if t.body == todo.body {
				t.status = t.status.next()
			}
		}

		m.groups = updateGroups(m.todosList)
		return m, func() tea.Msg {
			nextX := clamp(0, m.position.x+1, 2)
			m.position.x = nextX
			currGroup := m.activeGroup()
			for i, t := range currGroup.items {
				if t.body == todo.body {
					nextY := i

					return newPositionMsg(nextX, nextY)
				}
			}
			return newPositionMsg(0, 0)
		}

	case key.Matches(msg, keys.Help):
		m.help.ShowAll = !m.help.ShowAll
		return m, nil
	case key.Matches(msg, keys.Quit):
		m.todosList.writeToFile()
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
		x: x,
		y: y,
	}
}

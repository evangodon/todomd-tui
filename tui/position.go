package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evangodon/todomd/task"
)

type Position struct {
	Y int
	X int
}

func (p Position) GoRight() Position {
	p.X++
	return p
}

func (p Position) GoLeft() Position {
	p.X--
	return p
}

func (p Position) GoUp() Position {
	p.Y -= 1
	return p
}

func (p Position) GoDown() Position {
	p.Y += 1
	return p
}

type PositionMsg struct {
	nextPos Position
}

func newPositionMsg(p Position) tea.Cmd {
	return func() tea.Msg {
		return PositionMsg{
			nextPos: p,
		}
	}
}

func (m model) handleNewPosition(msg PositionMsg) (model, tea.Cmd) {
	nextPos := msg.nextPos

	// Validate X Position
	nextPos.X = task.Clamp(0, nextPos.X, 2)

	// Validate Y Position
	nextPos.X = task.Clamp(0, nextPos.X, 2)
	nextGroup := m.GetGroup(nextPos.X)
	maxY := len(nextGroup.Tasks()) - 1
	nextPos.Y = task.Clamp(0, m.position.Y, maxY)

	m.position = nextPos
	return m, nil
}

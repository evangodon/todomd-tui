package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evangodon/todo/internal"
)

func (m model) handleTextInputMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.textinput.enabled = false
			m.textinput.input.Blur()
			return m, nil
		case "enter":
			m.textinput.enabled = false
			m.textinput.input.Blur()
			todo := internal.NewTodo(m.textinput.input.Value(), internal.UncompletedStatus)
			m.textinput.input.Reset()
			m.todosList.AddTodo(todo)
			m.groups = updateGroups(m.todosList)
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.textinput.input, cmd = m.textinput.input.Update(msg)
	return m, cmd
}

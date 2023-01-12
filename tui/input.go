package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evangodon/todomd/task"
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
			todo := task.New(m.textinput.input.Value(), task.UncompletedStatus, m.textinput.parent)
			m.textinput.input.Reset()
			return m, NewTaskMsg(addTask, todo)
		}
	}
	var cmd tea.Cmd
	m.textinput.input, cmd = m.textinput.input.Update(msg)
	return m, cmd
}

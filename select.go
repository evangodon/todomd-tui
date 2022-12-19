package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type selectModel struct {
	todos     []*Todo
	position  int
	selection *Todo
	err       error
}

func initialSelectModel(todos []*Todo) selectModel {
	var selection *Todo
	if len(todos) > 0 {
		selection = todos[0]
	}

	return selectModel{
		todos:     todos,
		position:  0,
		selection: selection,
		err:       nil,
	}
}

func (m selectModel) Init() tea.Cmd {
	return nil
}

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			m.position = (m.position + 1) % len(m.todos)
			m.selection = m.todos[m.position]
			return m, nil
		case tea.KeyUp:
			m.position--
			if m.position < 0 {
				m.position = len(m.todos) - 1
			}
			m.selection = m.todos[m.position]
			return m, nil
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, cmd
}

func (m selectModel) View() string {
	doc := strings.Builder{}
	for i, todo := range m.todos {
		selected := i == m.position
		indicator := "  "
		if selected {
			indicator = "→ "
		}
		choice := fmt.Sprintf("%s%s", indicator, todo.render(40))
		doc.WriteString(choice)
	}

	doc.WriteString("    (esc to quit)")

	container := lipgloss.NewStyle().Padding(1, 2).Render

	out := container(doc.String())
	return out
}

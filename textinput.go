package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type textinputModel struct {
	textInput textinput.Model
	err       error
}

type (
	errMsg error
)

func initialModel() textinputModel {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()

	return textinputModel{
		textInput: ti,
		err:       nil,
	}
}

func (m textinputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m textinputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m textinputModel) View() string {
	return fmt.Sprintf(
		"Name of the todo:\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

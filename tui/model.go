package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evangodon/todo/internal"
)

type Window struct {
	Width  int
	Height int
}

type Position struct {
	Y int
	X int
}

type TextInput struct {
	input   textinput.Model
	enabled bool
}

type model struct {
	todosList *internal.Todos
	window    Window
	position  Position
	groups    internal.GroupsByStatus
	help      help.Model
	err       error
	textinput TextInput
}

func NewInteractiveModel(file string) model {
	list := internal.NewTodos(file)
	return model{
		todosList: list,
		window:    Window{Width: 0, Height: 0},
		position:  Position{Y: 0, X: 0},
		groups:    list.GroupByStatus(),
		help:      help.New(),
		err:       nil,
		textinput: TextInput{
			input:   textinput.New(),
			enabled: false,
		},
	}
}

func updateGroups(todos *internal.Todos) internal.GroupsByStatus {
	return todos.GroupByStatus()
}

type todofileRead struct{ err error }

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		if err := m.todosList.ParseFile(); err != nil {
			return todofileRead{err: err}
		}

		return todofileRead{}
	}
}

func (m model) activeGroup() internal.Group {
	group := map[int]internal.Group{
		0: m.groups.Uncompleted,
		1: m.groups.InProgress,
		2: m.groups.Completed,
	}

	return group[m.position.X]
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.textinput.enabled {
		return m.handleTextInputMsg(msg)
	}
	switch msg := msg.(type) {
	case todofileRead:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m.groups = m.todosList.GroupByStatus()
		return m, nil
	case tea.KeyMsg:
		m, cmd := m.handleKey(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		m.window.Height = msg.Height
		m.window.Width = msg.Width - appMargin
		m.help.Width = msg.Width
		return m, nil
	case Position:
		m.position = msg
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	groups := RenderGroups([]internal.Group{
		m.groups.Uncompleted,
		m.groups.InProgress,
		m.groups.Completed,
	}, m.window, m.position, m.textinput)

	helpbar := m.help.View(keys)

	height := m.window.Height - lipgloss.Height(helpbar)
	groups = container.Height(height).MaxHeight(height).Render(groups)

	s := strings.Builder{}
	s.WriteString(groups)
	s.WriteString("\n")
	s.WriteString(helpbar)

	return s.String()
}

const (
	appMargin = 1
)

var container = lipgloss.NewStyle().Margin(appMargin)

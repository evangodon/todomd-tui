package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evangodon/todomd/task"
	"github.com/evangodon/todomd/ui"
)

type Window struct {
	Width  int
	Height int
}

type TextInput struct {
	input   textinput.Model
	enabled bool
}

type model struct {
	todosList *task.List
	window    Window
	position  Position
	groups    task.GroupsByStatus
	help      help.Model
	err       error
	textinput TextInput
}

func NewModel(file string) model {
	list := task.NewList(file)
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

func updateGroups(todos *task.List) task.GroupsByStatus {
	return todos.GroupByStatus()
}

type fileReadMsg struct{ err error }

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		if err := m.todosList.ParseFile(); err != nil {
			return fileReadMsg{err: err}
		}

		return fileReadMsg{}
	}
}

func (m model) activeGroup() task.Group {
	group := map[int]task.Group{
		0: m.groups.Uncompleted,
		1: m.groups.InProgress,
		2: m.groups.Completed,
	}

	return group[m.position.X]
}

func (m model) GetGroup(xPos int) task.Group {
	group := map[int]task.Group{
		0: m.groups.Uncompleted,
		1: m.groups.InProgress,
		2: m.groups.Completed,
	}

	return group[xPos]
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.textinput.enabled {
		return m.handleTextInputMsg(msg)
	}

	switch msg := msg.(type) {
	case fileReadMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
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
	case PositionMsg:
		return m.handleNewPosition(msg)
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		header := ui.RedText.SetString("Error")
		return lipgloss.NewStyle().
			BorderForeground(ui.Color.Red).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Render(fmt.Sprintf("%s\n%s\n%s", header, m.err.Error(), "Press q to quit"))
	}
	groups := RenderGroups([]task.Group{
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

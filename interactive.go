package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

func (cmd Cmd) Interactive() *cli.Command {
	return &cli.Command{
		Name:    "interactive",
		Aliases: []string{"i"},
		Usage:   "Open in interactive view",
		Action: func(ctx *cli.Context) error {
			file := ctx.String("file")

			p := tea.NewProgram(newModel(file), tea.WithAltScreen())
			m, err := p.Run()
			if err != nil {
				return err
			}
			err = m.(model).err
			if err != nil {
				return err
			}

			return nil
		},
	}
}

type Window struct {
	width  int
	height int
}

type Position struct {
	y int
	x int
}

type model struct {
	todosList *Todos
	window    Window
	position  Position
	groups    groupsByStatus
	help      help.Model
	err       error
}

func newModel(file string) model {
	list := newTodos(file)
	return model{
		todosList: list,
		window: Window{
			width:  0,
			height: 0,
		},
		position: Position{
			y: 0,
			x: 0,
		},
		groups: list.groupByStatus(),
		help:   help.New(),
	}
}

func updateGroups(todos *Todos) groupsByStatus {
	return todos.groupByStatus()
}

type todofileRead struct{ err error }

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		if err := m.todosList.parseFile(); err != nil {
			return todofileRead{err: err}
		}

		return todofileRead{}
	}
}

func (m model) activeGroup() Group {
	group := map[int]Group{
		0: m.groups.uncompleted,
		1: m.groups.inProgress,
		2: m.groups.completed,
	}

	return group[m.position.x]
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case todofileRead:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m.groups = m.todosList.groupByStatus()
		return m, nil
	case tea.KeyMsg:
		m, cmd := m.handleKey(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		m.window.height = msg.Height
		m.window.width = msg.Width - appMargin
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
	groups := renderGroups([]Group{
		m.groups.uncompleted,
		m.groups.inProgress,
		m.groups.completed,
	}, termSize{
		width:  m.window.width,
		height: m.window.height,
	}, m.position)

	helpbar := m.help.View(keys)

	height := m.window.height - lipgloss.Height(helpbar)
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

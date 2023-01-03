package todoselect

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evangodon/todo/internal"
)

func New(todos []*internal.Todo) (*internal.Todo, error) {
	p := tea.NewProgram(initialModel(todos))

	m, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	return m.(model).selection, m.(model).err
}

type model struct {
	todos     []*internal.Todo
	position  int
	selection *internal.Todo
	err       error
	termWidth int
}

func initialModel(todos []*internal.Todo) model {
	var selection *internal.Todo
	if len(todos) > 0 {
		selection = todos[0]
	}

	return model{
		todos:     todos,
		position:  0,
		selection: selection,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		return m, nil
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
	}

	return m, cmd
}

func (m model) View() string {
	doc := strings.Builder{}
	for i, todo := range m.todos {
		selected := i == m.position
		indicator := "  "
		if selected {
			indicator = "→ "
		}
		choice := fmt.Sprintf("%s%s\n", indicator, todo.Render(m.termWidth-2, false))
		doc.WriteString(choice)
	}

	doc.WriteString("    (esc to quit)")

	container := lipgloss.NewStyle().Padding(1, 2).Render

	out := container(doc.String())
	return out
}

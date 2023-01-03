package todoselect

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evangodon/todomd/task"
)

func New(tasks []*task.Task) (*task.Task, error) {
	p := tea.NewProgram(initialModel(tasks))

	m, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	return m.(model).selection, m.(model).err
}

type model struct {
	tasks     []*task.Task
	position  int
	selection *task.Task
	err       error
	termWidth int
}

func initialModel(tasks []*task.Task) model {
	var selection *task.Task
	if len(tasks) > 0 {
		selection = tasks[0]
	}

	return model{
		tasks:     tasks,
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
			m.position = (m.position + 1) % len(m.tasks)
			m.selection = m.tasks[m.position]
			return m, nil
		case tea.KeyUp:
			m.position--
			if m.position < 0 {
				m.position = len(m.tasks) - 1
			}
			m.selection = m.tasks[m.position]
			return m, nil
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m model) View() string {
	doc := strings.Builder{}
	for i, task := range m.tasks {
		selected := i == m.position
		indicator := "  "
		if selected {
			indicator = "→ "
		}
		task.SetMaxWidth(m.termWidth - 2)
		choice := fmt.Sprintf("%s%s\n", indicator, task.Render())
		doc.WriteString(choice)
	}

	doc.WriteString("    (esc to quit)")

	container := lipgloss.NewStyle().Padding(1, 2).Render

	out := container(doc.String())
	return out
}

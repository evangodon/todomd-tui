package task

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/evangodon/todomd/ui"
)

// A list of todos of the same type.
// A group knows how to render itself in markdown or in the terminal
type Group struct {
	status   Status
	items    []Task
	maxWidth int
	selected int
}

func newGroup(status Status, items []Task) *Group {
	return &Group{
		status:   status,
		items:    items,
		maxWidth: 100,
		selected: -1,
	}
}

func (g Group) Items() []Task {
	return g.items
}

func (g *Group) SetSelected(selected int) {
	g.selected = selected
}

func (g *Group) SetMaxWidth(width int) {
	g.maxWidth = width
}

func (g Group) Status() Status {
	return g.status
}

func (g *Group) addTodo(todo Task) {
	g.items = append(g.items, todo)
}

func (g *Group) removeTodo(todo Task) error {
	n := len(g.items)
	index := sort.Search(n, func(i int) bool {
		return g.items[i].Body() == todo.Body()
	})

	if index == n {
		return errors.New("todo not found in TodosList")
	}

	g.items = append(g.items[:index], g.items[index+1:]...)
	return nil
}

var listContainer = lg.NewStyle().Padding(0, 1)

func (g *Group) Render() string {
	out := strings.Builder{}
	status := Status(g.status)
	header := formatHeader(status, g.selected >= 0)
	out.WriteString(header)
	out.WriteString("\n")

	if len(g.items) == 0 {
		out.WriteString(ui.DimText.Render("\n(none)"))
	}

	for i, t := range g.items {
		t.SetMaxWidth(g.maxWidth)
		t.SetIsSelected(i == g.selected)
		s := t.Render()

		out.WriteString("\n" + s)
	}

	return listContainer.Width(g.maxWidth).
		Render(out.String())
}

func (g Group) ToMarkdown() string {
	data := statusData[g.status]
	s := strings.Builder{}
	b := "#"
	var header string

	switch g.status {
	case UncompletedStatus:
		header = fmt.Sprintf("%s %s", b, data.header)
	case InProgressStatus:
		header = fmt.Sprintf("%s %s", b, data.header)
	case CompletedStatus:
		header = fmt.Sprintf("%s %s", b, data.header)
	default:
		return "UNKNOWN"
	}
	s.WriteString(header)
	s.WriteString("\n\n")
	if len(g.items) == 0 {
		return s.String()
	}

	for _, todo := range g.items {
		line := fmt.Sprintf("%s %s\n", statusData[g.status].mdIcon, todo.Body())
		s.WriteString(line)
	}
	s.WriteString("\n")

	return s.String()
}

func formatHeader(status Status, colActive bool) string {
	data := statusData[status]
	if colActive {
		data.header = ui.SelectedText.Render(data.header)
	}

	switch status {
	case UncompletedStatus:
		return fmt.Sprintf("%s %s", data.termHeaderBlock, data.header)
	case InProgressStatus:
		return fmt.Sprintf("%s %s", data.termHeaderBlock, data.header)
	case CompletedStatus:
		return fmt.Sprintf("%s %s", data.termHeaderBlock, data.header)
	default:
		return "UNKNOWN"
	}
}

func (g *Group) Width() int {
	return lipgloss.Width(g.Render())
}

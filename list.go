package main

import (
	"fmt"
	"strings"

	lg "github.com/charmbracelet/lipgloss"
)

// A list of todos of the same type
type List struct {
	status Status
	items  []*Todo
}

func newList(status Status) *List {
	return &List{
		status: status,
		items:  []*Todo{},
	}
}

func (l *List) addTodo(todo Todo) {
	l.items = append(l.items, &todo)
}

var listContainer = lg.NewStyle().Padding(1).Render

func (l *List) render() string {
	out := strings.Builder{}
	status := Status(l.status)
	out.WriteString(HeaderStyle(status))
	out.WriteString("\n\n")

	for _, t := range l.items {
		out.WriteString(t.render())
	}

	return listContainer(out.String())
}

func HeaderStyle(status Status) string {
	switch status {
	case todoStatus:
		b := lg.NewStyle().Background(blue).Render(block)
		return fmt.Sprintf("%s  %s", b, "TODO")
	case inProgressStatus:
		b := lg.NewStyle().Background(yellow).Render(block)
		return fmt.Sprintf("%s  %s", b, "IN-PROGRESS")
	case completedStatus:
		b := lg.NewStyle().Background(mauve).Render(block)
		return fmt.Sprintf("%s  %s", b, "COMPLETED")
	default:
		return "UNKNOWN"
	}
}

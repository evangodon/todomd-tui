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

func (l *List) String() string {
	data := statusData[l.status]
	b := "#"
	var header string
	out := strings.Builder{}
	switch l.status {
	case uncompletedStatus:
		header = fmt.Sprintf("%s %s", b, data.header)
	case inProgressStatus:
		header = fmt.Sprintf("%s %s", b, data.header)
	case completedStatus:
		header = fmt.Sprintf("%s %s", b, data.header)
	default:
		return "UNKNOWN"
	}
	out.WriteString(header)
	out.WriteString("\n\n")

	for _, todo := range l.items {
		line := fmt.Sprintf("%s %s\n", statusData[l.status].mdIcon, todo.body)
		out.WriteString(line)
	}
	out.WriteString("\n")

	return out.String()
}

func HeaderStyle(status Status) string {
	data := statusData[status]

	switch status {
	case uncompletedStatus:
		return fmt.Sprintf("%s %s", data.termHeaderBlock, data.header)
	case inProgressStatus:
		return fmt.Sprintf("%s %s", data.termHeaderBlock, data.header)
	case completedStatus:
		return fmt.Sprintf("%s %s", data.termHeaderBlock, data.header)
	default:
		return "UNKNOWN"
	}
}

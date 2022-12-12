package main

import (
	"fmt"
	"strings"

	lg "github.com/charmbracelet/lipgloss"
)

// A list of todos of the same type
type Group struct {
	status Status
	items  []*Todo
}

func newGroup(status Status) *Group {
	return &Group{
		status: status,
		items:  []*Todo{},
	}
}

func (g *Group) addTodo(todo Todo) {
	g.items = append(g.items, &todo)
}

var listContainer = lg.NewStyle().Padding(1).Render

func (g *Group) render() string {
	out := strings.Builder{}
	status := Status(g.status)
	out.WriteString(HeaderStyle(status))
	out.WriteString("\n\n")

	for _, t := range g.items {
		out.WriteString(t.render())
	}

	return listContainer(out.String())
}

func (g *Group) String() string {
	data := statusData[g.status]
	b := "#"
	var header string
	out := strings.Builder{}
	switch g.status {
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

	for _, todo := range g.items {
		line := fmt.Sprintf("%s %s\n", statusData[g.status].mdIcon, todo.body)
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

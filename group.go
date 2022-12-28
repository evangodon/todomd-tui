package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	lg "github.com/charmbracelet/lipgloss"
)

// A list of todos of the same type
type Group struct {
	status   Status
	items    []*Todo
	maxWidth int
}

func newGroup(status Status) *Group {
	return &Group{
		status:   status,
		items:    []*Todo{},
		maxWidth: 100,
	}
}

func (g *Group) addTodo(todo Todo) {
	g.items = append(g.items, &todo)
}

func (g *Group) removeTodo(todo Todo) error {
	n := len(g.items)
	index := sort.Search(n, func(i int) bool {
		return g.items[i].body == todo.body
	})

	if index == n {
		return errors.New("todo not found in TodosList")
	}

	g.items = append(g.items[:index], g.items[index+1:]...)
	return nil
}

var listContainer = lg.NewStyle().Padding(0, 1).Render

func (g *Group) render() string {
	out := strings.Builder{}
	status := Status(g.status)
	out.WriteString(HeaderStyle(status))
	out.WriteString("\n")

	if len(g.items) == 0 {
		out.WriteString(dimText("\n(none)"))
		return out.String()
	}

	for _, t := range g.items {
		s := t.render(g.maxWidth)
		out.WriteString("\n" + s)
	}

	return listContainer(out.String())
}

func (g *Group) String() string {
	data := statusData[g.status]
	s := strings.Builder{}
	b := "#"
	var header string

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
	s.WriteString(header)
	s.WriteString("\n\n")

	for _, todo := range g.items {
		line := fmt.Sprintf("%s %s\n", statusData[g.status].mdIcon, todo.body)
		s.WriteString(line)
	}
	s.WriteString("\n")

	return s.String()
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

func (g *Group) Width() int {
	return lipgloss.Width(g.render())
}

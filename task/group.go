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
	tasks    []Task
	maxWidth int
	selected int
}

func newGroup(status Status, tasks []Task) *Group {
	return &Group{
		status:   status,
		tasks:    tasks,
		maxWidth: 100,
		selected: -1,
	}
}

func (g Group) Tasks() []Task {
	return g.tasks
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

func (g *Group) addTask(task Task) {
	g.tasks = append(g.tasks, task)
}

func (g *Group) removeTask(task Task) error {
	n := len(g.tasks)
	index := sort.Search(n, func(i int) bool {
		return g.tasks[i].Body() == task.Body()
	})

	if index == n {
		return errors.New("todo not found in TodosList")
	}

	g.tasks = append(g.tasks[:index], g.tasks[index+1:]...)
	return nil
}

var listContainer = lg.NewStyle().Padding(0, 1)

func (g *Group) Render() string {
	out := strings.Builder{}
	status := Status(g.status)
	header := formatHeader(status, g.selected >= 0)
	out.WriteString(header)
	out.WriteString("\n")

	if len(g.tasks) == 0 {
		out.WriteString(ui.DimText.Render("\n(none)"))
	}

	tasks := g.Tasks()
	var lastParent Task
	for i, t := range tasks {
		if t.parent == nil {
			lastParent = t
		}
		t.SetMaxWidth(g.maxWidth)
		if i == g.selected {
			t.SetIsSelected(true)
		}

		out.WriteString("\n")
		if t.parent != nil && t.parent.Body() == lastParent.Body() {
			if i == len(tasks)-1 || tasks[i+1].parent == nil {
				out.WriteString("  └ " + t.Render())
				continue
			}

			out.WriteString("  ├ " + t.Render())
			continue
		}
		out.WriteString(t.Render())
	}

	return listContainer.Width(g.maxWidth).Render(out.String())
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
	if len(g.tasks) == 0 {
		return s.String()
	}

	for _, task := range g.tasks {
		if task.parent != nil {
			continue
		}
		line := fmt.Sprintf("%s %s\n", statusData[task.status].mdIcon, task.Body())
		s.WriteString(line)
		// TODO: figure out how to persist state of subtasks in file
		for _, subTask := range task.SubTasks() {
			s.WriteString(fmt.Sprintf("  %s %s\n", statusData[task.status].mdIcon, subTask.Body()))
		}
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

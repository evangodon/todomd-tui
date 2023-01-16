package task

import (
	"fmt"

	lg "github.com/charmbracelet/lipgloss"
	"github.com/evangodon/todomd/ui"
)

type Task struct {
	status     Status
	body       string
	maxWidth   int
	isSelected bool
	subTasks   []*Task
	parent     *Task
}

func New(body string, status Status, parent *Task) *Task {
	return &Task{
		status:   status,
		body:     body,
		subTasks: make([]*Task, 0),
		parent:   parent,
	}
}

func (t *Task) Status() Status {
	return t.status
}

func (t *Task) SetStatus(status Status) {
	t.status = status
	for _, sub := range t.subTasks {
		sub.SetStatus(status)
	}
	// TODO: if all subtasks have same  status, than change parent status to same
}

func (t *Task) SetMaxWidth(w int) {
	t.maxWidth = w
}

func (t *Task) SetIsSelected(s bool) {
	t.isSelected = s
}

func (t *Task) AddSubTask(s *Task) {
	t.subTasks = append(t.subTasks, s)
}

func (t *Task) RemoveSubTask(s *Task) {
	withRemoved := []*Task{}
	for _, t := range t.SubTasks() {
		if t.Body() != s.Body() {
			withRemoved = append(withRemoved, t)
		}
	}

	t.subTasks = withRemoved
}

func (t Task) Body() string {
	return t.body
}

func (t Task) SubTasks() []*Task {
	return t.subTasks
}

func (t Task) Parent() *Task {
	return t.parent
}

var (
	block = "▒▒"
	color = ui.Color
)

var statusData = map[Status]struct {
	termHeaderBlock string
	termIcon        string
	mdIcon          string
	header          string
}{
	UncompletedStatus: {
		termHeaderBlock: lg.NewStyle().Background(color.Blue).Render(block),
		termIcon:        "▢",
		mdIcon:          "- [ ]",
		header:          "TODO",
	},
	InProgressStatus: {
		termHeaderBlock: lg.NewStyle().Background(color.Yellow).Render(block),
		termIcon:        "◌",
		mdIcon:          "- [ ]",
		header:          "IN-PROGRESS",
	},
	CompletedStatus: {
		termHeaderBlock: lg.NewStyle().Background(color.Green).Render(block),
		termIcon:        "✓",
		mdIcon:          "- [x]",
		header:          "DONE",
	},
}

func (t Task) Render() string {
	var icon lg.Style
	data := statusData[t.Status()]
	body := lg.NewStyle().SetString(truncate(t.Body(), t.maxWidth-4))

	switch t.Status() {
	case UncompletedStatus:
		icon = ui.BlueText.SetString(data.termIcon)
	case InProgressStatus:
		icon = ui.YellowText.SetString(data.termIcon)
	case CompletedStatus:
		body = body.Faint(true)
		icon = ui.GreenText.SetString(data.termIcon)
	default:
		icon = lg.NewStyle()
	}

	if t.isSelected {
		icon = ui.RedText.Bold(true).SetString("ᐅ")
		body = body.Bold(true).Inline(true)
	}

	// subs := strings.Builder{}
	// moreRendered := false
	// for i, subTask := range t.SubTasks() {
	// 	subs.WriteString("\n")
	// 	if moreRendered || i < len(t.SubTasks())-1 {
	// 		subs.WriteString("  ├ " + subTask.Render())
	// 		continue
	// 	}
	// 	subs.WriteString("  └ " + subTask.Render())
	// 	moreRendered = true
	// }

	return fmt.Sprintf("%s %s", icon, body)
}

func (t Task) length() int {
	iconAndSpaces := 3
	return len(t.Body()) + iconAndSpaces
}

package task

import (
	"fmt"

	lg "github.com/charmbracelet/lipgloss"
	"github.com/evangodon/todomd/ui"
)

type Task struct {
	Status     Status
	body       string
	maxWidth   int
	isSelected bool
}

func New(body string, status Status) *Task {
	return &Task{
		Status: status,
		body:   body,
	}
}

func (t *Task) SetStatus(status Status) {
	t.Status = status
}

func (t *Task) SetMaxWidth(w int) {
	t.maxWidth = w
}

func (t *Task) SetIsSelected(s bool) {
	t.isSelected = s
}

func (t Task) Body() string {
	return t.body
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
	data := statusData[t.Status]
	body := lg.NewStyle().SetString(truncate(t.Body(), t.maxWidth-2))

	switch t.Status {
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
		body = body.Bold(true)
	}

	return lg.NewStyle().Inline(true).Render(fmt.Sprintf("%s %s", icon, body))
}

func (t Task) length() int {
	iconAndSpaces := 3
	return len(t.Body()) + iconAndSpaces
}

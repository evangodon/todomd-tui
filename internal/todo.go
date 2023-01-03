package internal

import (
	"fmt"

	lg "github.com/charmbracelet/lipgloss"
	"github.com/evangodon/todo/ui"
)

type Status int

const (
	UncompletedStatus Status = iota
	InProgressStatus
	CompletedStatus
)

func (s Status) String() string {
	switch s {
	case UncompletedStatus:
		return "TODO"
	case InProgressStatus:
		return "IN-PROGRESS"
	case CompletedStatus:
		return "DONE"
	default:
		return "UNKNOWN"
	}
}
func (s Status) Next() Status {
	return Status(Clamp(0, int(s)+1, 2))
}

func (s Status) Prev() Status {
	return Status(Clamp(0, int(s)-1, 2))
}

type Todo struct {
	Status Status
	body   string
}

func NewTodo(body string, status Status) *Todo {
	return &Todo{
		Status: status,
		body:   body,
	}
}

func (t *Todo) SetStatus(status Status) {
	t.Status = status
}

func (t Todo) Body() string {
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

func (t Todo) Render(maxwidth int, isSelected bool) string {
	var icon lg.Style
	data := statusData[t.Status]
	body := lg.NewStyle().SetString(truncate(t.Body(), maxwidth-4))

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

	if isSelected {
		icon = ui.RedText.Bold(true).SetString("ᐅ")
		body = body.Bold(true)
	}

	return lg.NewStyle().Inline(true).Render(fmt.Sprintf("%s %s", icon, body))
}

func (t Todo) length() int {
	iconAndSpaces := 3
	return len(t.Body()) + iconAndSpaces
}

package main

import (
	"fmt"

	lg "github.com/charmbracelet/lipgloss"
)

type Status int

const (
	uncompletedStatus Status = iota
	inProgressStatus
	completedStatus
)

func (s Status) String() string {
	switch s {
	case uncompletedStatus:
		return "TODO"
	case inProgressStatus:
		return "IN-PROGRESS"
	case completedStatus:
		return "DONE"
	default:
		return "UNKNOWN"
	}
}
func (s Status) next() Status {
	return Status(clamp(0, int(s)+1, 2))
}

func (s Status) prev() Status {
	return Status(clamp(0, int(s)-1, 2))
}

type Todo struct {
	status Status
	body   string
}

func newTodo(body string, status Status) *Todo {
	return &Todo{
		status: status,
		body:   body,
	}
}

var (
	block  = "▒▒"
	yellow = lg.Color("#f9e2af")
	blue   = lg.Color("#89b4fa")
	mauve  = lg.Color("#cba6f7")
	green  = lg.Color("#a6e3a1")
	red    = lg.Color("#f38ba8")
	white  = lg.Color("#ffffff")

	greenText    = lg.NewStyle().Foreground(green).Render
	redText      = lg.NewStyle().Foreground(red).Render
	blueText     = lg.NewStyle().Foreground(blue).Render
	yellowText   = lg.NewStyle().Foreground(yellow).Render
	selectedText = lg.NewStyle().Foreground(white).Bold(true).Render
	dimText      = lg.NewStyle().Faint(true).Render
)

var statusData = map[Status]struct {
	termHeaderBlock string
	termIcon        string
	mdIcon          string
	header          string
}{
	uncompletedStatus: {
		termHeaderBlock: lg.NewStyle().Background(blue).Render(block),
		termIcon:        "▢",
		mdIcon:          "- [ ]",
		header:          "TODO",
	},
	inProgressStatus: {
		termHeaderBlock: lg.NewStyle().Background(yellow).Render(block),
		termIcon:        "◌",
		mdIcon:          "- [ ]",
		header:          "IN-PROGRESS",
	},
	completedStatus: {
		termHeaderBlock: lg.NewStyle().Background(green).Render(block),
		termIcon:        "✓",
		mdIcon:          "- [x]",
		header:          "DONE",
	},
}

func (t Todo) render(maxwidth int, isSelected bool) string {
	var icon string
	data := statusData[t.status]
	body := truncate(t.body, maxwidth-4)

	switch t.status {
	case uncompletedStatus:
		icon = blueText(data.termIcon)
	case inProgressStatus:
		icon = yellowText(data.termIcon)
	case completedStatus:
		body = dimText(body)
		icon = greenText(data.termIcon)
	default:
		icon = ""
		body = t.body
	}

	if isSelected {
		icon = redText("ᐅ")
		body = selectedText(body)
	}

	return lg.NewStyle().Inline(true).Render(fmt.Sprintf("%s %s", icon, body))
}

func (t Todo) length() int {
	iconAndSpaces := 3
	return len(t.body) + iconAndSpaces
}

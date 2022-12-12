package main

import (
	"fmt"

	lg "github.com/charmbracelet/lipgloss"
)

type Status int

const (
	uncompletedStatus Status = iota
	completedStatus
	inProgressStatus
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

type Todo struct {
	status Status
	body   string
}

func newTodo(body string, status Status) Todo {
	return Todo{
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

	dimText    = lg.NewStyle().Faint(true).Render
	greenText  = lg.NewStyle().Foreground(green).Render
	redText    = lg.NewStyle().Foreground(red).Render
	blueText   = lg.NewStyle().Foreground(blue).Render
	yellowText = lg.NewStyle().Foreground(yellow).Render
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

func (t Todo) render() string {
	var icon string
	var body string
	data := statusData[t.status]

	switch t.status {
	case uncompletedStatus:
		body = t.body
		icon = blueText(data.termIcon)
	case inProgressStatus:
		body = t.body
		icon = yellowText(data.termIcon)
	case completedStatus:
		body = dimText(t.body)
		icon = greenText(data.termIcon)
	default:
		icon = ""
		body = t.body
	}

	return fmt.Sprintf("%s %s\n", icon, body)
}

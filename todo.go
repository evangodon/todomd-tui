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
)

var dimText = lg.NewStyle().Faint(true).Render
var greenText = lg.NewStyle().Foreground(green).Render

func (t Todo) render() string {
	var icon string
	var body string

	switch t.status {
	case uncompletedStatus:
		body = t.body
		icon = "▢"
	case inProgressStatus:
		body = t.body
		icon = "▢"
	case completedStatus:
		body = dimText(t.body)
		icon = "[✓]"
	default:
		icon = ""
		body = t.body
	}

	return fmt.Sprintf("%s %s\n", icon, body)
}

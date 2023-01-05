package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/evangodon/todomd/task"
)

const (
	gapWidth = 1
)

var (
	boldText = lipgloss.NewStyle().Bold(true).Render
)

func RenderGroups(groups []task.Group, win Window, pos Position, textinput TextInput) string {
	termWidth := win.Width
	totalNumGroups := len(groups)

	gap := strings.Repeat(" ", gapWidth)
	// Total width used up by gaps
	ttlGapSpace := len(gap) * (totalNumGroups - 1)
	ttlRenderedWidth := func() int {
		w := ttlGapSpace + (appMargin * 2)
		for _, g := range groups {
			w += g.Width()
		}
		return w
	}()

	parts := []string{}
	groupMaxWidth := termWidth / totalNumGroups
	over := ttlRenderedWidth >= termWidth
	if over {
		groupMaxWidth = termWidth/totalNumGroups - (ttlGapSpace+(appMargin*2))/(totalNumGroups-1)
	}
	for i, g := range groups {
		g.SetMaxWidth(groupMaxWidth)
		if i == pos.X {
			g.SetSelected(pos.Y)
		}

		spacebetween := gap
		if i%2 == 0 {
			spacebetween = ""
		}
		s := strings.Builder{}
		s.WriteString(g.Render())
		if g.Status() == task.UncompletedStatus && textinput.enabled {
			s.WriteString("\n")
			inputContainer := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
			out := inputContainer.Render(textinput.input.View())
			s.WriteString(out)
		}
		parts = append(parts, spacebetween, s.String(), spacebetween)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

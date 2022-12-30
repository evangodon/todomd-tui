package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type termSize struct {
	width  int
	height int
}

const (
	gapWidth = 1
)

var (
	boldText = lipgloss.NewStyle().Bold(true).Render
)

func renderGroups(groups []Group, t termSize, s Position) string {
	termWidth := t.width
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
		g.maxWidth = groupMaxWidth
		if i == s.x {
			g.selected = s.y
		}

		spacebetween := gap
		if i%2 == 0 {
			spacebetween = ""
		}
		parts = append(parts, spacebetween, g.render(), spacebetween)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

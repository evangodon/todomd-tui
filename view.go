package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

func (c Cmd) View() *cli.Command {
	return &cli.Command{
		Name:    "view",
		Aliases: []string{"v"},
		Usage:   "View all todos",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "width",
				Value: 0,
				Usage: "The max width",
			},
		},
		Action: func(ctx *cli.Context) error {
			file := ctx.String("file")
			definedWidth := ctx.Int("width")

			todosList := newTodos(file)
			if err := todosList.parseFile(); err != nil {
				return err
			}

			renderedGroups := []*Group{
				todosList.uncompleted,
				todosList.inProgress,
				todosList.completed,
			}

			termWidth, _ := c.TermSize()
			maxWidth := termWidth
			totalNumGroups := len(renderedGroups)
			if definedWidth > 0 {
				maxWidth = definedWidth
			}

			gap := strings.Repeat(" ", 3)
			// Total width used up by gaps
			ttlGapSpace := len(gap) * (totalNumGroups - 1)
			ttlRenderedWidth := func() int {
				w := ttlGapSpace
				for _, g := range renderedGroups {
					w += g.Width()
				}
				return w
			}()

			over := ttlRenderedWidth > maxWidth

			parts := []string{}
			groupMaxWidth := termWidth / totalNumGroups
			if over {
				groupMaxWidth = maxWidth/totalNumGroups - ttlGapSpace/(totalNumGroups-1)
			}
			for i, g := range renderedGroups {
				g.maxWidth = groupMaxWidth
				rendered := g.render()

				spacebetween := gap
				if i%2 == 0 {
					spacebetween = ""
				}
				parts = append(parts, spacebetween, rendered, spacebetween)
			}

			out := lipgloss.JoinHorizontal(lipgloss.Top, parts...)

			c.Log(logDefault, out)
			return nil
		},
	}
}

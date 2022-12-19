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
			if definedWidth > 0 {
				maxWidth = definedWidth
			}

			gap := strings.Repeat(" ", 8)
			ttlGapSpace := len(gap) * (len(renderedGroups) - 1)
			ttlRenderedWidth := func() int {
				w := ttlGapSpace
				for _, g := range renderedGroups {
					w += g.Width()
				}
				return w
			}()

			over := ttlRenderedWidth > maxWidth

			parts := []string{}
			for i, g := range renderedGroups {
				if over {
					maxWidth = maxWidth/len(renderedGroups) - ttlGapSpace/len(renderedGroups)
				}
				maxWidth = max(maxWidth, 10)
				g.maxWidth = maxWidth
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

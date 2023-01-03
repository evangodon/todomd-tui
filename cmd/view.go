package cmd

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/evangodon/todomd/task"
	"github.com/evangodon/todomd/tui"
	"github.com/evangodon/todomd/ui"
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

			todosList := task.NewList(file)
			if err := todosList.ParseFile(); err != nil {
				return err
			}
			println(lipgloss.HasDarkBackground())

			renderedGroups := []task.Group{
				todosList.CreateGroup(task.UncompletedStatus),
				todosList.CreateGroup(task.InProgressStatus),
				todosList.CreateGroup(task.CompletedStatus),
			}

			termWidth, termHeight := c.TermSize()
			out := tui.RenderGroups(
				renderedGroups,
				tui.Window{Width: termWidth, Height: termHeight},
				tui.Position{Y: -1, X: -1},
				tui.TextInput{},
			)

			ui.Log(ui.LogDefault, out)
			return nil
		},
	}
}

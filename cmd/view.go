package cmd

import (
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

			list := task.NewList(file)
			if err := list.ParseFile(); err != nil {
				return err
			}

			groupedByStatus := list.GroupByStatus()
			renderedGroups := []task.Group{
				groupedByStatus.Uncompleted,
				groupedByStatus.InProgress,
				groupedByStatus.Completed,
			}

			termWidth, termHeight := c.TermSize()
			out := tui.RenderGroups(
				renderedGroups,
				tui.Window{Width: termWidth, Height: termHeight},
				tui.Position{Y: -1, X: -1},
				tui.NewTaskInput{},
			)

			ui.Log(ui.LogDefault, out)
			return nil
		},
	}
}

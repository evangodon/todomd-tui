package cmd

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/evangodon/todo/internal"
	"github.com/evangodon/todo/tui"
	"github.com/evangodon/todo/ui"
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

			todosList := internal.NewTodos(file)
			if err := todosList.ParseFile(); err != nil {
				return err
			}
			println(lipgloss.HasDarkBackground())

			renderedGroups := []internal.Group{
				todosList.CreateGroup(internal.UncompletedStatus),
				todosList.CreateGroup(internal.InProgressStatus),
				todosList.CreateGroup(internal.CompletedStatus),
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

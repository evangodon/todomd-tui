package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evangodon/todomd/tui"
	"github.com/urfave/cli/v2"
)

func (cmd Cmd) Interactive() *cli.Command {
	return &cli.Command{
		Name:    "interactive",
		Aliases: []string{"i"},
		Usage:   "Open in interactive view",
		Action: func(ctx *cli.Context) error {
			file := ctx.String("file")

			p := tea.NewProgram(tui.NewModel(file), tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				return err
			}

			return nil
		},
	}
}

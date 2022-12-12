package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

func (c Cmd) View() *cli.Command {
	return &cli.Command{
		Name:    "view",
		Aliases: []string{"v"},
		Usage:   "View all todos",
		Action: func(ctx *cli.Context) error {
			file := ctx.String("file")

			todosList := newTodos(file)
			todosList.parseFile()

			gap := strings.Repeat(" ", 8)
			out := lipgloss.JoinHorizontal(
				lipgloss.Top,
				todosList.uncompleted.render(),
				gap,
				todosList.inProgress.render(),
				gap,
				todosList.completed.render(),
			)
			fmt.Print(out)
			return nil
		},
	}
}

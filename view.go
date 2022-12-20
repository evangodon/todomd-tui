package main

import (
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

			todosList := newTodos(file)
			if err := todosList.parseFile(); err != nil {
				return err
			}

			renderedGroups := []Group{
				todosList.createGroup(uncompletedStatus),
				todosList.createGroup(inProgressStatus),
				todosList.createGroup(completedStatus),
			}

			termWidth, termHeight := c.TermSize()
			out := renderGroups(
				renderedGroups,
				termSize{termWidth, termHeight},
				Position{y: -1, x: -1},
			)

			c.Log(logDefault, out)
			return nil
		},
	}
}

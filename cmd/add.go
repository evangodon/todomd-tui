package cmd

import (
	"github.com/evangodon/todomd/components/input"
	"github.com/evangodon/todomd/task"
	"github.com/evangodon/todomd/ui"
	"github.com/urfave/cli/v2"
)

func (cmd Cmd) Add() *cli.Command {
	return &cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "Add a todo",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Value: "",
				Usage: "The name of the todo",
			},
		},
		Action: func(ctx *cli.Context) error {
			filename := ctx.String("file")
			todosList := task.NewList(filename)
			if err := todosList.ParseFile(); err != nil {
				return err
			}

			todoName := input.New()

			if todoName == "" {
				return nil
			}

			todo := task.New(todoName, task.UncompletedStatus, nil)
			todosList.AddTask(todo)

			if err := todosList.WriteToFile(); err != nil {
				return err
			}
			ui.Log(ui.LogSuccess, "Added: "+todoName)

			return nil
		},
	}
}

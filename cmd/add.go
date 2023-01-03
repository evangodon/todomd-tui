package cmd

import (
	"github.com/evangodon/todo/components/input"
	"github.com/evangodon/todo/internal"
	"github.com/evangodon/todo/ui"
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
			todosList := internal.NewTodos(filename)
			if err := todosList.ParseFile(); err != nil {
				return err
			}

			todoName := input.New()

			if todoName == "" {
				return nil
			}

			todo := internal.NewTodo(todoName, internal.UncompletedStatus)
			todosList.AddTodo(todo)

			if err := todosList.WriteToFile(); err != nil {
				return err
			}
			ui.Log(ui.LogSuccess, "Added: "+todoName)

			return nil
		},
	}
}

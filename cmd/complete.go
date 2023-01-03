package cmd

import (
	"fmt"

	todoselect "github.com/evangodon/todomd/components/select"
	"github.com/evangodon/todomd/task"
	"github.com/evangodon/todomd/ui"
	"github.com/urfave/cli/v2"
)

func (cmd Cmd) Complete() *cli.Command {

	return &cli.Command{
		Name:    "complete",
		Aliases: []string{"c"},
		Usage:   "complete todo",
		Action: func(ctx *cli.Context) error {
			filename := ctx.String("file")
			todosList := task.NewList(filename)
			if err := todosList.ParseFile(); err != nil {
				return err
			}

			inProgressTodos := todosList.FilterByStatus(task.InProgressStatus)
			if len(inProgressTodos) == 0 {
				return cli.Exit("No todos in progress", 0)
			}

			todo, err := todoselect.New(inProgressTodos)
			if err != nil {
				return err
			}

			if todo == nil {
				return cli.Exit("", 0)
			}

			todo.SetStatus(task.CompletedStatus)
			todosList.WriteToFile()

			msg := fmt.Sprintf("completed '%s'", todo.Body())
			ui.Log(ui.LogSuccess, msg)
			return nil
		},
	}
}

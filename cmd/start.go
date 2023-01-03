package cmd

import (
	"fmt"

	todoselect "github.com/evangodon/todomd/components/select"
	"github.com/evangodon/todomd/task"
	"github.com/evangodon/todomd/ui"
	"github.com/urfave/cli/v2"
)

func (cmd Cmd) Start() *cli.Command {
	return &cli.Command{
		Name:    "start",
		Aliases: []string{"s"},
		Usage:   "start a todo",
		Action: func(ctx *cli.Context) error {
			filename := ctx.String("file")
			todosList := task.NewList(filename)
			if err := todosList.ParseFile(); err != nil {
				return err
			}

			uncompletedTodos := todosList.FilterByStatus(task.UncompletedStatus)
			if len(uncompletedTodos) == 0 {
				return cli.Exit("No unstarted todos", 0)
			}

			todo, err := todoselect.New(uncompletedTodos)
			if err != nil {
				return err
			}

			todo.Status = task.InProgressStatus

			todosList.WriteToFile()

			msg := fmt.Sprintf("started '%s'", todo.Body())
			ui.Log(ui.LogSuccess, msg)
			return nil
		},
	}
}

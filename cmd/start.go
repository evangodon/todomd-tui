package cmd

import (
	"fmt"

	todoselect "github.com/evangodon/todo/components/select"
	"github.com/evangodon/todo/internal"
	"github.com/evangodon/todo/ui"
	"github.com/urfave/cli/v2"
)

func (cmd Cmd) Start() *cli.Command {
	return &cli.Command{
		Name:    "start",
		Aliases: []string{"s"},
		Usage:   "start a todo",
		Action: func(ctx *cli.Context) error {
			filename := ctx.String("file")
			todosList := internal.NewTodos(filename)
			if err := todosList.ParseFile(); err != nil {
				return err
			}

			uncompletedTodos := todosList.FilterByStatus(internal.UncompletedStatus)
			if len(uncompletedTodos) == 0 {
				return cli.Exit("No unstarted todos", 0)
			}

			todo, err := todoselect.New(uncompletedTodos)
			if err != nil {
				return err
			}

			todo.Status = internal.InProgressStatus

			todosList.WriteToFile()

			msg := fmt.Sprintf("started '%s'", todo.Body())
			ui.Log(ui.LogSuccess, msg)
			return nil
		},
	}
}

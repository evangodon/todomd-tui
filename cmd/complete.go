package cmd

import (
	"fmt"

	todoselect "github.com/evangodon/todo/components/select"
	"github.com/evangodon/todo/internal"
	"github.com/evangodon/todo/ui"
	"github.com/urfave/cli/v2"
)

func (cmd Cmd) Complete() *cli.Command {

	return &cli.Command{
		Name:    "complete",
		Aliases: []string{"c"},
		Usage:   "complete todo",
		Action: func(ctx *cli.Context) error {
			filename := ctx.String("file")
			todosList := internal.NewTodos(filename)
			if err := todosList.ParseFile(); err != nil {
				return err
			}

			inProgressTodos := todosList.FilterByStatus(internal.InProgressStatus)
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

			todo.SetStatus(internal.CompletedStatus)
			todosList.WriteToFile()

			msg := fmt.Sprintf("completed '%s'", todo.Body())
			ui.Log(ui.LogSuccess, msg)
			return nil
		},
	}
}

package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
)

func (cmd Cmd) Complete() *cli.Command {

	return &cli.Command{
		Name:    "complete",
		Aliases: []string{"c"},
		Usage:   "complete todo",
		Action: func(ctx *cli.Context) error {
			filename := ctx.String("file")
			todosList := newTodos(filename)
			if err := todosList.parseFile(); err != nil {
				return err
			}

			inProgressTodos := todosList.inProgress.items
			if len(inProgressTodos) == 0 {
				return cli.Exit("No todos in progress", 0)
			}

			p := tea.NewProgram(initialSelectModel(inProgressTodos))

			m, err := p.Run()
			if err != nil {
				log.Fatal(err)
			}

			todo := m.(selectModel).selection

			if todo == nil {
				return cli.Exit("", 0)
			}

			todosList.completeTodo(*todo)
			todosList.writeToFile()

			msg := fmt.Sprintf("completed '%s'", todo.body)
			cmd.Log(logSuccess, msg)
			return nil
		},
	}
}

package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
)

func (cmd Cmd) Start() *cli.Command {
	return &cli.Command{
		Name:    "start",
		Aliases: []string{"s"},
		Usage:   "start a todo",
		Action: func(ctx *cli.Context) error {
			filename := ctx.String("file")
			todosList := newTodos(filename)
			if err := todosList.parseFile(); err != nil {
				return err
			}

			uncompletedTodos := todosList.uncompleted.items
			if len(uncompletedTodos) == 0 {
				return cli.Exit("No unstarted todos", 0)
			}
			p := tea.NewProgram(initialSelectModel(uncompletedTodos))

			m, err := p.Run()
			if err != nil {
				log.Fatal(err)
			}

			todo := m.(selectModel).selection

			todosList.startTodo(*todo)
			todosList.writeToFile()

			msg := fmt.Sprintf("started '%s'", todo.body)
			cmd.Log(logSuccess, msg)
			return nil
		},
	}
}

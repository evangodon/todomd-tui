package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
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
			todosList := newTodos(filename)
			if err := todosList.parseFile(); err != nil {
				return err
			}

			p := tea.NewProgram(initialTextinputModel())

			m, err := p.Run()
			if err != nil {
				log.Fatal(err)
			}

			todoName := m.(textinputModel).textInput.Value()

			if todoName == "" {
				return nil
			}

			todosList.uncompleted.addTodo(newTodo(todoName, uncompletedStatus))

			err = todosList.writeToFile()
			if err != nil {
				return err
			}
			cmd.Log(logSuccess, "Added: "+todoName)

			return nil
		},
	}
}

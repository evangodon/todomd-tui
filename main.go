package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	cmd := NewCmd()

	commands := []*cli.Command{
		cmd.Add(),
		cmd.View(),
	}

	app := &cli.App{
		Name:     "todo",
		Usage:    "create a todo",
		Commands: commands,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file",
				Value: "todo.md",
				Usage: "Loads todos from `FILE`",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

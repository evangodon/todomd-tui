package main

import (
	"os"

	"github.com/evangodon/todomd/cmd"
	"github.com/evangodon/todomd/ui"
	"github.com/urfave/cli/v2"
)

func main() {
	cmd := cmd.New()

	commands := []*cli.Command{
		cmd.Add(),
		cmd.View(),
		cmd.Start(),
		cmd.Complete(),
		cmd.Interactive(),
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
		ui.Log(ui.LogError, err.Error())
		os.Exit(1)
	}
}

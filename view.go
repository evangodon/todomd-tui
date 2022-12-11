package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

func (c Cmd) View() *cli.Command {
	return &cli.Command{
		Name:    "view",
		Aliases: []string{"v"},
		Usage:   "View all todos",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file",
				Value: "todo.md",
				Usage: "The file to read",
			},
		},
		Action: func(ctx *cli.Context) error {
			file := ctx.String("file")
			f, _ := os.Open(file)
			scanner := bufio.NewScanner(f)

			todos := make([]Todo, 0)
			var currentStatus Status

			for scanner.Scan() {
				line := scanner.Text()
				line = strings.TrimSpace(line)

				if line == "# TODO" {
					currentStatus = todoStatus
				}

				if line == "# IN-PROGRESS" {
					currentStatus = inProgressStatus
				}

				if line == "# DONE" {
					currentStatus = completedStatus
				}

				statusPattern := regexp.MustCompile(`- \[(x| )\]`)
				if matched := statusPattern.MatchString(line); matched {
					body := line[6:]
					todos = append(todos, newTodo(body, currentStatus))
				}
			}

			uncompletedList := newList(todoStatus)
			inProgressList := newList(inProgressStatus)
			completedList := newList(completedStatus)

			for _, todo := range todos {
				switch todo.status {
				case todoStatus:
					uncompletedList.addTodo(todo)
				case inProgressStatus:
					inProgressList.addTodo(todo)
				case completedStatus:
					completedList.addTodo(todo)
				default:
				}
			}

			gap := strings.Repeat(" ", 8)
			out := lipgloss.JoinHorizontal(
				lipgloss.Top,
				uncompletedList.render(),
				gap,
				inProgressList.render(),
				gap,
				completedList.render(),
			)
			fmt.Print(out)
			return nil
		},
	}
}

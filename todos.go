package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

type Todos struct {
	uncompleted *List
	inProgress  *List
	completed   *List
	filename    string
}

func newTodos(filename string) *Todos {

	todos := &Todos{
		uncompleted: newList(uncompletedStatus),
		inProgress:  newList(inProgressStatus),
		completed:   newList(completedStatus),
		filename:    filename,
	}

	return todos
}

func (td *Todos) parseFile() {
	f, err := os.OpenFile(td.filename, os.O_RDWR, 0777)
	if err != nil {
		log.Fatal("Error openeing file: ", err)
	}
	scanner := bufio.NewScanner(f)
	todos := make([]Todo, 0)
	var currentStatus Status

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.Contains(line, "# TODO") {
			currentStatus = uncompletedStatus
		}

		if strings.Contains(line, "# IN-PROGRESS") {
			currentStatus = inProgressStatus
		}

		if strings.Contains(line, "# DONE") {
			currentStatus = completedStatus
		}

		statusPattern := regexp.MustCompile(`- \[(x| )\]`)
		if matched := statusPattern.MatchString(line); matched {
			body := line[6:]
			body = strings.TrimSpace(body)
			todos = append(todos, newTodo(body, currentStatus))
		}
	}

	for _, todo := range todos {
		switch todo.status {
		case uncompletedStatus:
			td.uncompleted.addTodo(todo)
		case inProgressStatus:
			td.inProgress.addTodo(todo)
		case completedStatus:
			td.completed.addTodo(todo)
		default:
		}
	}
}

func (td *Todos) writeToFile() {
	update := strings.Builder{}

	_, err := update.WriteString(td.uncompleted.String())
	if err != nil {
		log.Fatal(err)
	}
	_, err = update.WriteString(td.inProgress.String())
	if err != nil {
		log.Fatal(err)
	}
	_, err = update.WriteString(td.completed.String())
	if err != nil {
		log.Fatal(err)
	}

	if err = os.Truncate(td.filename, 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	f, err := os.OpenFile(td.filename, os.O_RDWR, 0777)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(update.String())
	if err != nil {
		log.Fatal(err)
	}
}

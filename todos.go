package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Todos struct {
	uncompleted *Group
	inProgress  *Group
	completed   *Group
	filename    string
}

func newTodos(filename string) *Todos {
	return &Todos{
		uncompleted: newGroup(uncompletedStatus),
		inProgress:  newGroup(inProgressStatus),
		completed:   newGroup(completedStatus),
		filename:    filename,
	}
}

func (td *Todos) parseFile() error {
	f, err := os.OpenFile(td.filename, os.O_RDWR, 0777)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	scanner := bufio.NewScanner(f)
	todos := make([]Todo, 0)
	var currentStatus Status

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.Contains(line, statusData[uncompletedStatus].header) {
			currentStatus = uncompletedStatus
		}

		if strings.Contains(line, statusData[inProgressStatus].header) {
			currentStatus = inProgressStatus
		}

		if strings.Contains(line, statusData[completedStatus].header) {
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
	return nil
}

func (td *Todos) writeToFile() error {
	update := strings.Builder{}

	if _, err := update.WriteString(td.uncompleted.String()); err != nil {
		return err
	}
	if _, err := update.WriteString(td.inProgress.String()); err != nil {
		return err
	}
	if _, err := update.WriteString(td.completed.String()); err != nil {
		return err
	}

	if err := os.Truncate(td.filename, 0); err != nil {
		return fmt.Errorf("failed to truncate: %v", err)
	}

	f, err := os.OpenFile(td.filename, os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	_, err = f.WriteString(update.String())
	if err != nil {
		return err
	}
	return nil
}

func (td *Todos) startTodo(todo Todo) error {
	err := td.uncompleted.removeTodo(todo)
	if err != nil {
		return err
	}
	todo.status = inProgressStatus
	td.inProgress.addTodo(todo)

	return nil
}

func (td *Todos) completeTodo(todo Todo) error {
	err := td.inProgress.removeTodo(todo)
	if err != nil {
		return err
	}
	todo.status = completedStatus
	td.completed.addTodo(todo)

	return nil
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Todos struct {
	items    []*Todo
	filename string
}

func newTodos(filename string) *Todos {
	return &Todos{
		items:    make([]*Todo, 0),
		filename: filename,
	}
}

func (t *Todos) addTodo(todo *Todo) {
	t.items = append(t.items, todo)
}

func (t *Todos) filterByStatus(status Status) []*Todo {
	items := make([]*Todo, 0)
	for _, todo := range t.items {
		if todo.status == status {
			items = append(items, todo)
		}
	}

	return items
}

func (t *Todos) createGroup(status Status) Group {
	items := make([]Todo, 0)
	for _, todo := range t.items {
		if todo.status == status {
			items = append(items, *todo)
		}
	}

	return *newGroup(status, items)
}

type groupsByStatus struct {
	uncompleted Group
	inProgress  Group
	completed   Group
}

func (t *Todos) groupByStatus() groupsByStatus {
	groups := groupsByStatus{
		uncompleted: *newGroup(uncompletedStatus, []Todo{}),
		inProgress:  *newGroup(inProgressStatus, []Todo{}),
		completed:   *newGroup(completedStatus, []Todo{}),
	}

	for _, todo := range t.items {
		switch todo.status {
		case uncompletedStatus:
			groups.uncompleted.addTodo(*todo)
		case inProgressStatus:
			groups.inProgress.addTodo(*todo)
		case completedStatus:
			groups.completed.addTodo(*todo)
		}
	}

	return groups
}

func (td *Todos) parseFile() error {
	f, err := os.OpenFile(td.filename, os.O_RDWR, 0777)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	scanner := bufio.NewScanner(f)
	todos := make([]*Todo, 0)
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

	td.items = todos

	return nil
}

func (td *Todos) writeToFile() error {
	update := strings.Builder{}
	groupsByStatus := td.groupByStatus()

	if _, err := update.WriteString(groupsByStatus.uncompleted.String()); err != nil {
		return err
	}
	if _, err := update.WriteString(groupsByStatus.inProgress.String()); err != nil {
		return err
	}
	if _, err := update.WriteString(groupsByStatus.completed.String()); err != nil {
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

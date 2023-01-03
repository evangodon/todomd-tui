package task

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type List struct {
	tasks    []*Task
	filename string
}

func NewList(filename string) *List {
	return &List{
		tasks:    make([]*Task, 0),
		filename: filename,
	}
}

func (t List) Tasks() []*Task {
	return t.tasks
}

func (t *List) AddTask(task *Task) {
	t.tasks = append(t.tasks, task)
}

func (t *List) FilterByStatus(status Status) []*Task {
	items := make([]*Task, 0)
	for _, task := range t.tasks {
		if task.Status == status {
			items = append(items, task)
		}
	}

	return items
}

func (t *List) CreateGroup(status Status) Group {
	items := make([]Task, 0)
	for _, task := range t.tasks {
		if task.Status == status {
			items = append(items, *task)
		}
	}

	return *newGroup(status, items)
}

type GroupsByStatus struct {
	Uncompleted Group
	InProgress  Group
	Completed   Group
}

func (t *List) GroupByStatus() GroupsByStatus {
	groups := GroupsByStatus{
		Uncompleted: *newGroup(UncompletedStatus, []Task{}),
		InProgress:  *newGroup(InProgressStatus, []Task{}),
		Completed:   *newGroup(CompletedStatus, []Task{}),
	}

	for _, task := range t.tasks {
		switch task.Status {
		case UncompletedStatus:
			groups.Uncompleted.addTask(*task)
		case InProgressStatus:
			groups.InProgress.addTask(*task)
		case CompletedStatus:
			groups.Completed.addTask(*task)
		}
	}

	return groups
}

func (td *List) ParseFile() error {
	f, err := os.OpenFile(td.filename, os.O_RDWR, 0777)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	scanner := bufio.NewScanner(f)
	tasks := make([]*Task, 0)
	var currentStatus Status

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.Contains(line, statusData[UncompletedStatus].header) {
			currentStatus = UncompletedStatus
		}

		if strings.Contains(line, statusData[InProgressStatus].header) {
			currentStatus = InProgressStatus
		}

		if strings.Contains(line, statusData[CompletedStatus].header) {
			currentStatus = CompletedStatus
		}

		statusPattern := regexp.MustCompile(`- \[(x| )\]`)
		if matched := statusPattern.MatchString(line); matched {
			body := line[6:]
			body = strings.TrimSpace(body)
			tasks = append(tasks, New(body, currentStatus))
		}
	}

	td.tasks = tasks

	return nil
}

func (td *List) WriteToFile() error {
	update := strings.Builder{}
	groupsByStatus := td.GroupByStatus()

	if _, err := update.WriteString(groupsByStatus.Uncompleted.ToMarkdown()); err != nil {
		return err
	}
	if _, err := update.WriteString(groupsByStatus.InProgress.ToMarkdown()); err != nil {
		return err
	}
	if _, err := update.WriteString(groupsByStatus.Completed.ToMarkdown()); err != nil {
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

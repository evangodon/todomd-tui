package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList_ParseFile(t *testing.T) {
	list := NewList("../testdata/test-todo.md")
	err := list.ParseFile()
	assert.NoError(t, err)

	assert.Equal(t, 18, len(list.Tasks()))

	groups := list.GroupByStatus()

	assert.Equal(t, 7, len(groups.Uncompleted.tasks))
	assert.Equal(t, 1, len(groups.InProgress.tasks))
	assert.Equal(t, 10, len(groups.Completed.tasks))
}

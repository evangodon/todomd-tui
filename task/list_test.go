package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList_ParseFile(t *testing.T) {
	list := NewList("../testdata/test-todo.md")
	err := list.ParseFile()
	assert.NoError(t, err)

	assert.Equal(t, 16, len(list.items))

	groups := list.GroupByStatus()

	assert.Equal(t, 5, len(groups.Uncompleted.items))
	assert.Equal(t, 1, len(groups.InProgress.items))
	assert.Equal(t, 10, len(groups.Completed.items))
}

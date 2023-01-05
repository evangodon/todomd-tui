package task

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestTask_Render(t *testing.T) {

	tests := []struct {
		task     Task
		maxWidth int
	}{
		{
			task:     *New("Validate config file and all modules", UncompletedStatus, nil),
			maxWidth: 20,
		},
		{
			task:     *New("use xdg paths for config", UncompletedStatus, nil),
			maxWidth: 10,
		},
	}

	for _, tt := range tests {
		tt.task.SetMaxWidth(tt.maxWidth)
		out := tt.task.Render()
		assert.GreaterOrEqual(t, tt.maxWidth, lipgloss.Width(out))
	}

}

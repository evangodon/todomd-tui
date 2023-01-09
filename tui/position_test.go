package tui

import (
	"testing"
)

func Test_model_handleNewPosition(t *testing.T) {
	tests := []struct {
		name    string
		nextPos Position
		want    Position
	}{
		{
			name: "1",
			nextPos: Position{
				Y: 0,
				X: -1,
			},
			want: Position{
				Y: 0,
				X: 0,
			},
		},
		{
			name: "2",
			nextPos: Position{
				Y: 0,
				X: 5,
			},
			want: Position{
				Y: 0,
				X: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel("todo.md")

			msg := PositionMsg{
				nextPos: tt.nextPos,
			}
			got, _ := m.handleNewPosition(msg)
			if got.position.X != tt.want.X {
				t.Errorf("model.handleTextInputMsg() got = %v, want %v", got.position, tt.want)
			}
		})
	}
}

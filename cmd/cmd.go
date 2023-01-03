package cmd

import (
	"os"

	"golang.org/x/term"
)

type Cmd struct{}

func New() *Cmd {
	return &Cmd{}
}

func (c Cmd) TermSize() (width int, height int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 150, 150
	}

	return width, height
}

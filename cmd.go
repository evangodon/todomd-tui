package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type Cmd struct{}

func NewCmd() *Cmd {
	return &Cmd{}
}

type logType int

const (
	logSuccess logType = iota
	logError
	logInfo
	logDefault
)

func (c Cmd) TermSize() (width int, height int) {

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 150, 150
	}

	return width, height
}

func (c Cmd) Log(logtype logType, msg string) {
	var s string
	switch logtype {
	case logSuccess:
		icon := greenText("✓")
		s = fmt.Sprintf("%s %s", icon, msg)
	case logError:
		icon := redText("✘")
		s = fmt.Sprintf("%s %s", icon, msg)
	case logInfo:
		icon := "i"
		s = fmt.Sprintf("%s %s", icon, msg)
	case logDefault:
		s = msg
	}

	fmt.Println("\n" + s)
}

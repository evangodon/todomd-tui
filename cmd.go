package main

import "fmt"

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

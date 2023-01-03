package ui

import (
	"fmt"
)

type logType int

const (
	LogSuccess logType = iota
	LogError
	LogInfo
	LogDefault
)

func Log(logtype logType, msg string) {
	var s string
	switch logtype {
	case LogSuccess:
		icon := GreenText.SetString("✓")
		s = fmt.Sprintf("%s %s", icon, msg)
	case LogError:
		icon := RedText.SetString("✘")
		s = fmt.Sprintf("%s %s", icon, msg)
	case LogInfo:
		icon := "i"
		s = fmt.Sprintf("%s %s", icon, msg)
	case LogDefault:
		s = msg
	}

	fmt.Println(s)
}

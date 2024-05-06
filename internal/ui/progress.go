package ui

import (
	"io"
)

type ProgressGroup interface {
	Add(message string) Progress
	Start() error
	Stop() error
}

type NewProgressGroupType func(output io.Writer) ProgressGroup

func NewProgressGroup(output io.Writer) ProgressGroup {
	if StdoutIsTerminal() {
		return newTerminalProgressGroup(output)
	} else {
		return newScriptProgressGroup(output)
	}
}

type Progress interface {
	Start()
	SetCurrent(value int)
	SetSuccess()
	SetError()
}

func NewProgress(output io.Writer, message string) Progress {
	if StdoutIsTerminal() {
		return newTerminalProgress(output, message)
	} else {
		return newScriptProgress(output, message)
	}
}

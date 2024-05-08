package ui

import (
	"io"
)

type ProgressGroup interface {
	Add(message string, resources string) Progress
	Start() error
	Stop() error
}

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

func NewProgress(output io.Writer, message string, resources string) Progress {
	if StdoutIsTerminal() {
		return newTerminalProgress(output, message, resources)
	} else {
		return newScriptProgress(output, message, resources)
	}
}

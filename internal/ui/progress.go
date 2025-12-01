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
	if StdoutIsTerminal() && TerminalWidth() >= 80 {
		return newTerminalProgressGroup(output)
	}
	return newScriptProgressGroup(output)
}

type Progress interface {
	Start()
	SetCurrent(value int)
	SetSuccess()
	SetError()
}

func NewProgress(output io.Writer, message string, resources string) Progress {
	if StdoutIsTerminal() && TerminalWidth() > 80 {
		return newTerminalProgress(output, message, resources)
	}
	return newScriptProgress(output, message, resources)
}

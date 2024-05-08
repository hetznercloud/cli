package ui

import (
	"fmt"
	"io"
)

type scriptProgressGroup struct {
	output   io.Writer
	progress []Progress
}

func newScriptProgressGroup(output io.Writer) *scriptProgressGroup {
	return &scriptProgressGroup{output: output}
}

func (p *scriptProgressGroup) Add(message string, resources string) Progress {
	progress := newScriptProgress(p.output, message, resources)
	p.progress = append(p.progress, progress)
	return progress

}

func (p *scriptProgressGroup) Start() error {
	for _, progress := range p.progress {
		progress.Start()
	}
	return nil
}

func (p *scriptProgressGroup) Stop() error {
	return nil
}

type scriptProgress struct {
	output    io.Writer
	message   string
	resources string
}

func newScriptProgress(output io.Writer, message string, resources string) *scriptProgress {
	return &scriptProgress{output: output, message: message, resources: resources}
}

func (p *scriptProgress) print(status string) {
	result := p.message
	if p.resources != "" {
		result += fmt.Sprintf(" %s", p.resources)
	}
	result += " ..."
	if status != "" {
		result += fmt.Sprintf(" %s", status)
	}
	fmt.Fprintln(p.output, result)
}

func (p *scriptProgress) Start() {
	p.print("")
}

func (p *scriptProgress) SetCurrent(value int) {
}

func (p *scriptProgress) SetSuccess() {
	p.print("done")
}

func (p *scriptProgress) SetError() {
	p.print("failed")
}

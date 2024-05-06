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

func (p *scriptProgressGroup) Add(message string) Progress {
	progress := newScriptProgress(p.output, message)
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
	output  io.Writer
	message string
}

func newScriptProgress(output io.Writer, message string) *scriptProgress {
	return &scriptProgress{output: output, message: message}
}

func (p *scriptProgress) Start() {
	fmt.Fprintf(p.output, "%s ...\n", p.message)
}

func (p *scriptProgress) SetCurrent(value int) {}

func (p *scriptProgress) SetSuccess() {
	fmt.Fprintf(p.output, "%s ... done\n", p.message)
}

func (p *scriptProgress) SetError() {
	fmt.Fprintf(p.output, "%s ... failed\n", p.message)
}

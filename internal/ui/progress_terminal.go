package ui

import (
	"fmt"
	"io"

	"github.com/cheggaaa/pb/v3"
)

type terminalProgressGroup struct {
	output io.Writer
	el     *pb.Pool
}

func newTerminalProgressGroup(output io.Writer) *terminalProgressGroup {
	return &terminalProgressGroup{output: output, el: pb.NewPool()}
}

func (p *terminalProgressGroup) Add(message string, resources string) Progress {
	progress := newTerminalProgress(p.output, message, resources)
	p.el.Add(progress.el)
	return progress
}

func (p *terminalProgressGroup) Start() error {
	return p.el.Start()
}

func (p *terminalProgressGroup) Stop() error {
	return p.el.Stop()
}

const (
	terminalProgressTpl = `{{ string . "message" }} {{ percent . "%3.f%%" | blue }} {{ etime . | blue }} {{ string . "resources" | blue }}`

	terminalProgressRunningTpl = ` {{ cycle . "⠋" "⠙" "⠹" "⠸" "⠼" "⠴" "⠦" "⠧" "⠇" "⠏" | blue }} ` + terminalProgressTpl
	terminalProgressSuccessTpl = ` {{ green "✓" }} ` + terminalProgressTpl
	terminalProgressErrorTpl   = ` {{ red "✗" }} ` + terminalProgressTpl
)

type terminalProgress struct {
	el *pb.ProgressBar
}

func newTerminalProgress(output io.Writer, message string, resources string) *terminalProgress {
	p := &terminalProgress{pb.New(100)}
	p.el.SetWriter(output)
	p.el.SetTemplateString(terminalProgressRunningTpl)
	p.el.Set("message", fmt.Sprintf("%-60s", message))
	p.el.Set("resources", resources)
	return p
}

func (p *terminalProgress) Start() {
	p.el.Start()
}

func (p *terminalProgress) SetCurrent(value int) {
	p.el.SetCurrent(int64(value))
}

func (p *terminalProgress) SetSuccess() {
	p.el.SetCurrent(int64(100))
	p.el.SetTemplateString(terminalProgressSuccessTpl)
	p.el.Finish()
}

func (p *terminalProgress) SetError() {
	p.el.SetCurrent(100)
	p.el.SetTemplate(terminalProgressErrorTpl)
	p.el.Finish()
}

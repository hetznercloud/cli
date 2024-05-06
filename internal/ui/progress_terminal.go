package ui

import (
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

func (p *terminalProgressGroup) Add(message string) Progress {
	progress := newTerminalProgress(p.output, message)
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
	terminalProgressRunningTpl = `{{ cycle . "⠋" "⠙" "⠹" "⠸" "⠼" "⠴" "⠦" "⠧" "⠇" "⠏" }} {{ string . "message" }} {{ percent . }}`
	terminalProgressSuccessTpl = `{{ green "✓" }} {{ string . "message" }} {{ percent . }}`
	terminalProgressErrorTpl   = `{{ red "✗" }} {{ string . "message" }} {{ percent . }}`
)

type terminalProgress struct {
	el *pb.ProgressBar
}

func newTerminalProgress(output io.Writer, message string) *terminalProgress {
	p := &terminalProgress{pb.New(100)}
	p.el.SetWriter(output)
	p.el.SetTemplateString(terminalProgressRunningTpl)
	p.el.Set("message", message)
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

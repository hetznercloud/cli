package terminal

import (
	"golang.org/x/term"

	"github.com/hetznercloud/cli/internal/ui"
)

//go:generate go run go.uber.org/mock/mockgen -package terminal -destination zz_terminal_mock.go . Terminal

type Terminal interface {
	StdoutIsTerminal() bool
	ReadPassword(fd int) ([]byte, error)
}

type DefaultTerminal struct{}

func (DefaultTerminal) StdoutIsTerminal() bool {
	return ui.StdoutIsTerminal()
}

func (DefaultTerminal) ReadPassword(fd int) ([]byte, error) {
	return term.ReadPassword(fd)
}

var _ Terminal = DefaultTerminal{}

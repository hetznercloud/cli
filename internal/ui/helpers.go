package ui

import (
	"os"

	"golang.org/x/term"
)

// StdoutIsTerminal returns whether the CLI is run in a terminal.
func StdoutIsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

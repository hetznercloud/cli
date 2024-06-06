package ui

import (
	"os"

	"golang.org/x/term"
)

// StdoutIsTerminal returns whether the CLI is run in a terminal.
func StdoutIsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// TerminalWidth returns the width of the terminal.
func TerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0
	}

	return width
}

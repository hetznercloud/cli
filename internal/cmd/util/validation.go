package util

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// posArgRegex matches the placeholders for positional arguments in the usage string of a cobra.Command.
// Breakdown:
// 1. `\[.*?]` matches everything in square brackets, but does not capture them, so it is ignored.
// 2. `\(.*?\)` same as (1) but with parentheses.
// 3. `--.*?<.*?>` matches flags and ignores them, so that `--flag <arg>` is not counted as a positional argument.
// 4. `<([a-zA-Z0-9-_]*?)>\.\.\.` matches variadic arguments and captures the name of the argument in group 2.
// 5. `<([a-zA-Z0-9-_]*?)>` matches regular arguments and captures the name of the argument in group 1.
var posArgRegex = regexp.MustCompile(`\[.*?]|\(.*?\)|--.*?<.*?>|<([a-zA-Z0-9-_]*?)>\.\.\.|<([a-zA-Z0-9-_]*?)>`)

// Validate checks if the number of positional arguments matches the usage string of a cobra.Command.
// It returns an error if the number of arguments does not match the usage string (if it is not variadic) or if an
// argument is missing.
// If your command usage has optional positional arguments, use ValidateLenient instead.
func Validate(cmd *cobra.Command, args []string) error {
	return validate(cmd, args, false)
}

// ValidateLenient checks if the number of positional arguments matches the usage string of a cobra.Command.
// In contrast to Validate, it does not return an error if there are more arguments than expected.
// This can be useful for commands that have optional positional arguments.
func ValidateLenient(cmd *cobra.Command, args []string) error {
	return validate(cmd, args, true)
}

func validate(cmd *cobra.Command, args []string, lenient bool) error {
	use := cmd.Use
	paren := cmd.Parent()
	for paren != nil {
		use = paren.Use + " " + use
		paren = paren.Parent()
	}

	matches := posArgRegex.FindAllStringSubmatch(cmd.Use, -1)
	isVariadic := false
	var expected []string
	for _, match := range matches {
		if match[1] != "" {
			expected = append(expected, match[1])
			isVariadic = true
		} else if match[2] != "" {
			expected = append(expected, match[2])
		}
	}
	if len(args) > len(expected) && !isVariadic && !lenient {
		_, _ = fmt.Fprintln(os.Stderr, use)
		_, _ = fmt.Fprintln(os.Stderr, strings.Repeat(" ", len(use)+1)+"^")
		return fmt.Errorf("expected exactly %d positional argument(s), but got %d", len(expected), len(args))
	}

	for i := 0; i < len(expected); i++ {
		if i >= len(args) {
			idx := strings.Index(use, "<"+expected[i]+">")
			if idx != -1 {
				_, _ = fmt.Fprintln(os.Stderr, use)
				_, _ = fmt.Fprintln(os.Stderr, strings.Repeat(" ", idx+1)+strings.Repeat("^", len(expected[i])))
			}
			return fmt.Errorf("expected argument(s) %s at position %d", expected[i], i+1)
		}
	}
	return nil
}

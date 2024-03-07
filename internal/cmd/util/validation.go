package util

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var posArgRegex = regexp.MustCompile(`\[.*?]|\(.*?\)|--.*?<.*?>|<(.*?)>`) // matches all positional args in group 1

func Validate(cmd *cobra.Command, args []string) error {
	matches := posArgRegex.FindAllStringSubmatch(cmd.Use, -1)
	var expected []string
	for _, match := range matches {
		if match[1] != "" {
			expected = append(expected, match[1])
		}
	}
	if len(args) > len(expected) {
		return fmt.Errorf("expected exactly %d positional arguments, but got %d", len(expected), len(args))
	}
	for i := 0; i < len(expected); i++ {
		if i >= len(args) {
			idx := strings.Index(cmd.Use, "<"+expected[i]+">")
			if idx != -1 {
				_, _ = fmt.Fprintln(os.Stderr, cmd.Use)
				_, _ = fmt.Fprintln(os.Stderr, strings.Repeat(" ", idx+1)+strings.Repeat("^", len(expected[i])))
			}
			return fmt.Errorf("expected argument %s at position %d", strings.ReplaceAll(expected[i], "-", " "), i+1)
		}
	}
	return nil
}

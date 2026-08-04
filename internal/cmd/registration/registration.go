package registration

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const errorAnnotation = "hcloud.dev/command-construction-errors"

func Record(cmd *cobra.Command, err error) {
	if err == nil {
		return
	}
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	message := err.Error()
	if previous := cmd.Annotations[errorAnnotation]; previous != "" {
		message = previous + "\n" + message
	}
	cmd.Annotations[errorAnnotation] = message
}

func Error(root *cobra.Command) error {
	var errs []error
	var visit func(*cobra.Command)
	visit = func(cmd *cobra.Command) {
		if message := cmd.Annotations[errorAnnotation]; message != "" {
			for line := range strings.Lines(message) {
				errs = append(errs, fmt.Errorf("%s: %s", cmd.CommandPath(), strings.TrimSpace(line)))
			}
		}
		for _, child := range cmd.Commands() {
			visit(child)
		}
	}
	visit(root)
	return errors.Join(errs...)
}

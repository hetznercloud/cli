package util

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/registration"
)

func MarkFlagRequired(cmd *cobra.Command, flag string) {
	if err := cmd.MarkFlagRequired(flag); err != nil {
		registration.Record(cmd, fmt.Errorf("mark flag %q required: %w", flag, err))
	}
}

func MarkFlagFilename(cmd *cobra.Command, flag string, extensions ...string) {
	if err := cmd.MarkFlagFilename(flag, extensions...); err != nil {
		registration.Record(cmd, fmt.Errorf("mark flag %q as filename: %w", flag, err))
	}
}

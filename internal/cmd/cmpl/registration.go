package cmpl

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/registration"
)

func RegisterFlagCompletion(cmd *cobra.Command, flag string, completion cobra.CompletionFunc) {
	if err := cmd.RegisterFlagCompletionFunc(flag, completion); err != nil {
		registration.Record(cmd, fmt.Errorf("completion for flag %q: %w", flag, err))
	}
}

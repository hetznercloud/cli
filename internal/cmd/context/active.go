package context

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/state"
)

func newActiveCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "active [FLAGS]",
		Short:                 "Show active context",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  state.Wrap(s, runActive),
	}
	return cmd
}

func runActive(s state.State, cmd *cobra.Command, _ []string) error {
	if os.Getenv("HCLOUD_TOKEN") != "" {
		_, _ = fmt.Fprintln(os.Stderr, "Warning: HCLOUD_TOKEN is set. The active context will have no effect.")
	}
	if cfg := s.Config(); cfg.ActiveContext != nil {
		cmd.Println(cfg.ActiveContext.Name)
	}
	return nil
}

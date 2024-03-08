package context

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func newActiveCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "active",
		Short:                 "Show active context",
		Args:                  util.ValidateExact,
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
	if ctx := s.Config().ActiveContext(); ctx != nil {
		cmd.Println(ctx.Name)
	}
	return nil
}

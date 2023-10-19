package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/state"
)

func newActiveCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "active [FLAGS]",
		Short:                 "Show active context",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.Wrap(runActive),
	}
	return cmd
}

func runActive(cli *state.State, cmd *cobra.Command, args []string) error {
	if cli.Config.ActiveContext != nil {
		fmt.Println(cli.Config.ActiveContext.Name)
	}
	return nil
}

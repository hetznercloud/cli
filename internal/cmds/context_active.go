package cmds

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func newContextActiveCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "active [FLAGS]",
		Short:                 "Show active context",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.Wrap(runContextActive),
	}
	return cmd
}

func runContextActive(cli *state.State, cmd *cobra.Command, args []string) error {
	if cli.Config.ActiveContext != nil {
		fmt.Println(cli.Config.ActiveContext.Name)
	}
	return nil
}

package context

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "context [FLAGS]",
		Short:                 "Manage contexts",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newCreateCommand(cli),
		newActiveCommand(cli),
		newUseCommand(cli),
		newDeleteCommand(cli),
		newListCommand(cli),
	)
	return cmd
}

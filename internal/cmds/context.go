package cmds

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewContextCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "context [FLAGS]",
		Short:                 "Manage contexts",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newContextCreateCommand(cli),
		newContextActiveCommand(cli),
		newContextUseCommand(cli),
		newContextDeleteCommand(cli),
		newContextListCommand(cli),
	)
	return cmd
}

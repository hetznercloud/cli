package cmds

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewLocationCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "location",
		Short:                 "Manage locations",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newLocationListCommand(cli),
		newLocationDescribeCommand(cli),
	)
	return cmd
}

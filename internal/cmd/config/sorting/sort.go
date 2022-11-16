package sorting

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewSortCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "sort COMMAND",
		Short:                 "Configure the default sorting order for a command",
		Args:                  cobra.MinimumNArgs(2),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(
		newSetCommand(cli),
		newListCommand(cli),
	)

	return cmd
}

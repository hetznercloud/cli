package cmds

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewISOCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "iso",
		Short:                 "Manage ISOs",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newISOListCommand(cli),
		newISODescribeCommand(cli),
	)
	return cmd
}

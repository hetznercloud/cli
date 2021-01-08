package cmds

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewDatacenterCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "datacenter",
		Short:                 "Manage datacenters",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newDatacenterListCommand(cli),
		newDatacenterDescribeCommand(cli),
	)
	return cmd
}

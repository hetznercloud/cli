package datacenter

import (
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "datacenter",
		Short:                 "Manage datacenters",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		listCmd.CobraCommand(cli.Context, hcapi2.NewClient(cli.Client()), cli),
		describeCmd.CobraCommand(cli.Context, hcapi2.NewClient(cli.Client()), cli),
	)
	return cmd
}

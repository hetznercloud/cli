package loadbalancertype

import (
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCommand(cli *state.State, client hcapi2.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "load-balancer-type",
		Short:                 "Manage Load Balancer types",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		describeCmd.CobraCommand(cli.Context, client, cli),
		ListCmd.CobraCommand(cli.Context, client, cli, cli.Config.SubcommandDefaults[cmd.Use]),
	)
	return cmd
}

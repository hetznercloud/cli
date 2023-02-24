package loadbalancer

import (
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCommand(cli *state.State, client hcapi2.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "load-balancer",
		Short:                 "Manage Load Balancers",
		Aliases:               []string{"loadbalancer"},
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		CreateCommand.CobraCommand(cli.Context, client, cli, cli),
		ListCmd.CobraCommand(cli.Context, client, cli),
		DescribeCmd.CobraCommand(cli.Context, client, cli),
		deleteCmd.CobraCommand(cli.Context, client, cli),
		updateCmd.CobraCommand(cli.Context, client, cli),
		labelCmds.AddCobraCommand(cli.Context, client, cli),
		labelCmds.RemoveCobraCommand(cli.Context, client, cli),
		newAddTargetCommand(cli),
		newRemoveTargetCommand(cli),
		newChangeAlgorithmCommand(cli),
		newUpdateServiceCommand(cli),
		newDeleteServiceCommand(cli),
		newAddServiceCommand(cli),
		newEnableProtectionCommand(cli),
		newDisableProtectionCommand(cli),
		AttachToNetworkCommand.CobraCommand(cli.Context, client, cli, cli),
		DetachFromNetworkCommand.CobraCommand(cli.Context, client, cli, cli),
		newEnablePublicInterfaceCommand(cli),
		newDisablePublicInterfaceCommand(cli),
		newChangeTypeCommand(cli),
		newMetricsCommand(cli),
		setRDNSCmd.CobraCommand(cli.Context, client, cli, cli),
	)
	return cmd
}

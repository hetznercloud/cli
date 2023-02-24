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
		AddTargetCommand.CobraCommand(cli.Context, client, cli, cli),
		RemoveTargetCommand.CobraCommand(cli.Context, client, cli, cli),
		ChangeAlgorithmCommand.CobraCommand(cli.Context, client, cli, cli),
		UpdateServiceCommand.CobraCommand(cli.Context, client, cli, cli),
		DeleteServiceCommand.CobraCommand(cli.Context, client, cli, cli),
		AddServiceCommand.CobraCommand(cli.Context, client, cli, cli),
		EnableProtectionCommand.CobraCommand(cli.Context, client, cli, cli),
		DisableProtectionCommand.CobraCommand(cli.Context, client, cli, cli),
		AttachToNetworkCommand.CobraCommand(cli.Context, client, cli, cli),
		DetachFromNetworkCommand.CobraCommand(cli.Context, client, cli, cli),
		EnablePublicInterfaceCommand.CobraCommand(cli.Context, client, cli, cli),
		DisablePublicInterfaceCommand.CobraCommand(cli.Context, client, cli, cli),
		ChangeTypeCommand.CobraCommand(cli.Context, client, cli, cli),
		MetricsCommand.CobraCommand(cli.Context, client, cli, cli),
		setRDNSCmd.CobraCommand(cli.Context, client, cli, cli),
	)
	return cmd
}

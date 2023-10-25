package loadbalancer

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
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
		CreateCmd.CobraCommand(cli.Context, client, cli, cli),
		ListCmd.CobraCommand(cli.Context, client, cli),
		DescribeCmd.CobraCommand(cli.Context, client, cli),
		DeleteCmd.CobraCommand(cli.Context, client, cli, cli),
		UpdateCmd.CobraCommand(cli.Context, client, cli),
		LabelCmds.AddCobraCommand(cli.Context, client, cli),
		LabelCmds.RemoveCobraCommand(cli.Context, client, cli),
		AddTargetCmd.CobraCommand(cli.Context, client, cli, cli),
		RemoveTargetCmd.CobraCommand(cli.Context, client, cli, cli),
		ChangeAlgorithmCmd.CobraCommand(cli.Context, client, cli, cli),
		UpdateServiceCmd.CobraCommand(cli.Context, client, cli, cli),
		DeleteServiceCmd.CobraCommand(cli.Context, client, cli, cli),
		AddServiceCmd.CobraCommand(cli.Context, client, cli, cli),
		EnableProtectionCmd.CobraCommand(cli.Context, client, cli, cli),
		DisableProtectionCmd.CobraCommand(cli.Context, client, cli, cli),
		AttachToNetworkCmd.CobraCommand(cli.Context, client, cli, cli),
		DetachFromNetworkCmd.CobraCommand(cli.Context, client, cli, cli),
		EnablePublicInterfaceCmd.CobraCommand(cli.Context, client, cli, cli),
		DisablePublicInterfaceCmd.CobraCommand(cli.Context, client, cli, cli),
		ChangeTypeCmd.CobraCommand(cli.Context, client, cli, cli),
		MetricsCmd.CobraCommand(cli.Context, client, cli, cli),
		SetRDNSCmd.CobraCommand(cli.Context, client, cli, cli),
	)
	return cmd
}

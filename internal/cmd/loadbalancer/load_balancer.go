package loadbalancer

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "load-balancer",
		Short:                 "Manage Load Balancers",
		Aliases:               []string{"loadbalancer"},
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		CreateCmd.CobraCommand(s),
		ListCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		UpdateCmd.CobraCommand(s),
		LabelCmds.AddCobraCommand(s),
		LabelCmds.RemoveCobraCommand(s),
		AddTargetCmd.CobraCommand(s),
		RemoveTargetCmd.CobraCommand(s),
		ChangeAlgorithmCmd.CobraCommand(s),
		UpdateServiceCmd.CobraCommand(s),
		DeleteServiceCmd.CobraCommand(s),
		AddServiceCmd.CobraCommand(s),
		EnableProtectionCmd.CobraCommand(s),
		DisableProtectionCmd.CobraCommand(s),
		AttachToNetworkCmd.CobraCommand(s),
		DetachFromNetworkCmd.CobraCommand(s),
		EnablePublicInterfaceCmd.CobraCommand(s),
		DisablePublicInterfaceCmd.CobraCommand(s),
		ChangeTypeCmd.CobraCommand(s),
		MetricsCmd.CobraCommand(s),
		SetRDNSCmd.CobraCommand(s),
	)
	return cmd
}

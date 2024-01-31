package loadbalancer

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
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

	util.AddGroup(cmd, "general", "General",
		ListCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		UpdateCmd.CobraCommand(s),
		LabelCmds.AddCobraCommand(s),
		LabelCmds.RemoveCobraCommand(s),
		ChangeAlgorithmCmd.CobraCommand(s),
		ChangeTypeCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "protection", "Protection",
		EnableProtectionCmd.CobraCommand(s),
		DisableProtectionCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "target", "Targets",
		AddTargetCmd.CobraCommand(s),
		RemoveTargetCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "service", "Services",
		AddServiceCmd.CobraCommand(s),
		UpdateServiceCmd.CobraCommand(s),
		DeleteServiceCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "network", "Network",
		AttachToNetworkCmd.CobraCommand(s),
		DetachFromNetworkCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "public-interface", "Public Interface",
		EnablePublicInterfaceCmd.CobraCommand(s),
		DisablePublicInterfaceCmd.CobraCommand(s),
	)

	cmd.AddCommand(
		MetricsCmd.CobraCommand(s),
		SetRDNSCmd.CobraCommand(s),
	)
	return cmd
}

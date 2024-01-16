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
	cmd.AddGroup(&cobra.Group{ID: "general", Title: "General"})
	cmd.AddGroup(&cobra.Group{ID: "protection", Title: "Protection"})
	cmd.AddGroup(&cobra.Group{ID: "target", Title: "Targets"})
	cmd.AddGroup(&cobra.Group{ID: "service", Title: "Services"})
	cmd.AddGroup(&cobra.Group{ID: "network", Title: "Network"})
	cmd.AddGroup(&cobra.Group{ID: "public-interface", Title: "Public Interface"})

	cmd.AddCommand(
		util.WithGroup("general", ListCmd.CobraCommand(s)),
		util.WithGroup("general", DescribeCmd.CobraCommand(s)),
		util.WithGroup("general", DeleteCmd.CobraCommand(s)),
		util.WithGroup("general", UpdateCmd.CobraCommand(s)),
		util.WithGroup("general", LabelCmds.AddCobraCommand(s)),
		util.WithGroup("general", LabelCmds.RemoveCobraCommand(s)),
		util.WithGroup("general", ChangeAlgorithmCmd.CobraCommand(s)),
		util.WithGroup("general", ChangeTypeCmd.CobraCommand(s)),

		util.WithGroup("protection", EnableProtectionCmd.CobraCommand(s)),
		util.WithGroup("protection", DisableProtectionCmd.CobraCommand(s)),

		util.WithGroup("target", AddTargetCmd.CobraCommand(s)),
		util.WithGroup("target", RemoveTargetCmd.CobraCommand(s)),

		util.WithGroup("service", AddServiceCmd.CobraCommand(s)),
		util.WithGroup("service", UpdateServiceCmd.CobraCommand(s)),
		util.WithGroup("service", DeleteServiceCmd.CobraCommand(s)),

		util.WithGroup("network", AttachToNetworkCmd.CobraCommand(s)),
		util.WithGroup("network", DetachFromNetworkCmd.CobraCommand(s)),

		util.WithGroup("public-interface", EnablePublicInterfaceCmd.CobraCommand(s)),
		util.WithGroup("public-interface", DisablePublicInterfaceCmd.CobraCommand(s)),

		util.WithGroup("", MetricsCmd.CobraCommand(s)),
		util.WithGroup("", SetRDNSCmd.CobraCommand(s)),
	)
	return cmd
}

package zone

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/cmd/zone/rrset"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "zone",
		Aliases:               []string{"dns", "zones"},
		Short:                 "Manage DNS Zones and Zone RRSets (records)",
		Long:                  "For more details, see the documentation for Zones https://docs.hetzner.cloud/reference/cloud#zones or Zone RRSets https://docs.hetzner.cloud/reference/cloud#zone-rrsets.",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}

	util.AddGroup(cmd, "general", "General",
		DescribeCmd.CobraCommand(s),
		CreateCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		ListCmd.CobraCommand(s),
		LabelCmds.AddCobraCommand(s),
		LabelCmds.RemoveCobraCommand(s),
		ChangeTTLCmd.CobraCommand(s),
		ChangePrimaryNameserversCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "protection", "Protection",
		ChangeProtectionCmds.EnableCobraCommand(s),
		ChangeProtectionCmds.DisableCobraCommand(s),
	)

	util.AddGroup(cmd, "zonefile", "BIND Zone file",
		ExportZonefileCmd.CobraCommand(s),
		ImportZonefileCmd.CobraCommand(s),
	)

	util.AddGroup(cmd, "rrset", "Resource Record Sets (RRSets)",
		rrset.NewCommand(s),
		// Aliases for simple RRSet commands
		rrset.SetRecordsCmd.CobraCommand(s),
		rrset.AddRecordsCmd.CobraCommand(s),
		rrset.RemoveRecordsCmd.CobraCommand(s),
	)

	return cmd
}

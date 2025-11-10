package rrset

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DisableProtectionCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:                   "disable-protection <zone> <name> <type> change",
			Args:                  util.ValidateLenient,
			Short:                 "Disable resource protection for a Zone RRSet",
			ValidArgsFunction:     cmpl.SuggestArgs(rrsetArgumentsCompletionFuncs(client)...),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		zoneIDOrName, rrsetName, rrsetType := args[0], args[1], args[2]

		zoneIDOrName, err := util.ParseZoneIDOrName(zoneIDOrName)
		if err != nil {
			return fmt.Errorf("failed to convert Zone name to ascii: %w", err)
		}

		zone, _, err := s.Client().Zone().Get(s, zoneIDOrName)
		if err != nil {
			return err
		}
		if zone == nil {
			return fmt.Errorf("Zone not found: %s", zoneIDOrName)
		}

		rrset, _, err := s.Client().Zone().GetRRSetByNameAndType(s, zone, rrsetName, hcloud.ZoneRRSetType(rrsetType))
		if err != nil {
			return err
		}
		if rrset == nil {
			return fmt.Errorf("Zone RRSet not found: %s %s", rrsetName, rrsetType)
		}

		opts, err := getChangeProtectionOpts(false, args[3:])
		if err != nil {
			return err
		}

		return changeProtection(s, cmd, rrset, false, opts)
	},
}

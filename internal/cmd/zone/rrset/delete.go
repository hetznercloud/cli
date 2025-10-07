package rrset

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Unlike other delete commands, this does not support deleting multiple entities at once.
// As the rrset has three identifiers that need to be passed, its hard to find a good UX for deleting
// multiple at the same time.

var DeleteCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "delete <zone> <name> <type>",
			Short:                 "Delete a Zone RRSet",
			Args:                  util.Validate,
			ValidArgsFunction:     cmpl.SuggestArgs(rrsetArgumentsCompletionFuncs(client)...),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		return cmd
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
			return fmt.Errorf("Zone RRSet not found: %s %s %s", zoneIDOrName, rrsetName, rrsetType)
		}

		result, _, err := s.Client().Zone().DeleteRRSet(s, rrset)
		if err != nil {
			return err
		}

		err = s.WaitForActions(s, cmd, result.Action)
		if err != nil {
			return err
		}

		cmd.Printf("Zone RRSet %s %s deleted\n", rrsetName, rrsetType)

		return nil
	},
	Experimental: experimental.DNS,
}

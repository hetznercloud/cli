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

var ChangeTTLCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "change-ttl (--ttl <ttl> | --unset) <zone> <name> <type>",
			Short:                 "Changes the Time To Live (TTL) of a Zone RRSet",
			ValidArgsFunction:     cmpl.SuggestArgs(rrsetArgumentsCompletionFuncs(client)...),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().Int("ttl", 0, "Time To Live (TTL) of the Zone RRSet (required)")

		cmd.Flags().Bool("unset", false, "Unset the Time To Live of Zone RRSet (use the Zone default TTL instead)")

		cmd.MarkFlagsOneRequired("ttl", "unset")
		cmd.MarkFlagsMutuallyExclusive("ttl", "unset")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		zoneIDOrName, rrsetName, rrsetType := args[0], args[1], args[2]

		ttl, _ := cmd.Flags().GetInt("ttl")
		unset, _ := cmd.Flags().GetBool("unset")

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

		var opts hcloud.ZoneRRSetChangeTTLOpts
		if !unset {
			opts.TTL = &ttl
		}

		action, _, err := s.Client().Zone().ChangeRRSetTTL(s, rrset, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Changed TTL on Zone RRSet %s %s\n", rrset.Name, rrset.Type)
		return nil
	},
	Experimental: experimental.DNS,
}

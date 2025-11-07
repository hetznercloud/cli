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

var RemoveRecordsCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "remove-records (--record <value>... | --records-file <file>) <zone> <name> <type>",
			Short: "Remove records from a Zone RRSet.",
			Long: `Remove records from a Zone RRSet.

If the Zone RRSet doesn't contain any records, it will automatically be deleted.

` + recordsFileExample,
			ValidArgsFunction:     cmpl.SuggestArgs(rrsetArgumentsCompletionFuncs(client)...),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		addRecordsFlags(cmd)

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		zoneIDOrName, rrsetName, rrsetType := args[0], args[1], args[2]

		zoneIDOrName, err := util.ParseZoneIDOrName(zoneIDOrName)
		if err != nil {
			return fmt.Errorf("failed to convert Zone name to ascii: %w", err)
		}

		zone := &hcloud.Zone{Name: zoneIDOrName}

		rrset, _, err := s.Client().Zone().GetRRSetByNameAndType(s, zone, rrsetName, hcloud.ZoneRRSetType(rrsetType))
		if err != nil {
			return err
		}
		if rrset == nil {
			return fmt.Errorf("Zone RRSet not found: %s %s", rrsetName, rrsetType)
		}

		var opts hcloud.ZoneRRSetRemoveRecordsOpts

		opts.Records, err = recordsFromFlags(cmd)
		if err != nil {
			return err
		}

		// TXT: Format record values to simplify its usage
		if rrset.Type == hcloud.ZoneRRSetTypeTXT {
			FormatTXTRecords(cmd, opts.Records)
		}

		action, _, err := s.Client().Zone().RemoveRRSetRecords(s, rrset, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Removed records from Zone RRSet %s %s\n", rrset.Name, rrset.Type)
		return nil
	},
}

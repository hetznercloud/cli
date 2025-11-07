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

var AddRecordsCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "add-records (--record <value>... | --records-file <file>) <zone> <name> <type>",
			Short: "Add records to a Zone RRSet",
			Long: `Add records to a Zone RRSet.

If the Zone RRSet doesn't exist, it will automatically be created.

` + recordsFileExample,
			ValidArgsFunction:     cmpl.SuggestArgs(rrsetArgumentsCompletionFuncs(client)...),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		addRecordsFlags(cmd)

		cmd.Flags().Int("ttl", 0, "Time To Live (TTL) of the Zone RRSet")

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

		// This does not verify the rrset against the API, as the API will implicitly create a new RRSet when adding
		// records and the RRSet does not yet exist.
		rrset := &hcloud.ZoneRRSet{
			Zone: zone,
			Name: rrsetName,
			Type: hcloud.ZoneRRSetType(rrsetType),
		}

		var opts hcloud.ZoneRRSetAddRecordsOpts

		opts.Records, err = recordsFromFlags(cmd)
		if err != nil {
			return err
		}

		// TXT: Format record values to simplify its usage
		if rrset.Type == hcloud.ZoneRRSetTypeTXT {
			FormatTXTRecords(cmd, opts.Records)
		}

		if ttl, _ := cmd.Flags().GetInt("ttl"); ttl != 0 {
			opts.TTL = &ttl
		}

		action, _, err := s.Client().Zone().AddRRSetRecords(s, rrset, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Added records on Zone RRSet %s %s\n", rrset.Name, rrset.Type)
		return nil
	},
}

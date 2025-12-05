package rrset

import (
	"context"

	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/exp/zoneutil"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rrset",
		Aliases: []string{"record", "records"},
		Short:   "Manage Zone RRSets (records)",
		Long: `
For more details, see the documentation for Zones https://docs.hetzner.cloud/reference/cloud#zones
or Zone RRSets https://docs.hetzner.cloud/reference/cloud#zone-rrsets.

TXT records format must consist of one or many quoted strings of 255 characters. If the
user provider TXT records are not quoted, they will be formatted for you.`,
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}

	util.AddGroup(cmd, "general", "General",
		CreateCmd.CobraCommand(s),
		ListCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		ChangeTTLCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
		LabelCmds.AddCobraCommand(s),
		LabelCmds.RemoveCobraCommand(s),
	)

	util.AddGroup(cmd, "protection", "Protection",
		ChangeProtectionCmds.EnableCobraCommand(s),
		ChangeProtectionCmds.DisableCobraCommand(s),
	)

	util.AddGroup(cmd, "records", "Manage Records",
		SetRecordsCmd.CobraCommand(s),
		AddRecordsCmd.CobraCommand(s),
		RemoveRecordsCmd.CobraCommand(s),
	)

	return cmd
}

// addRecordsFlags adds flags for setting records. Used in [CreateCmd], [AddRecordsCmd] [RemoveRecordsCmd].
// To get the records use [recordsFromFlags].
func addRecordsFlags(cmd *cobra.Command) {
	cmd.Flags().String("records-file", "", "JSON file containing the records (conflicts with --record)")
	cmd.Flags().StringArray("record", []string{}, "List of records (can be specified multiple times, conflicts with --records-file)")
	cmd.MarkFlagsOneRequired("record", "records-file")
	cmd.MarkFlagsMutuallyExclusive("record", "records-file")
}

// recordsFromFlags parses the [hcloud.ZoneRRSetRecord] from `--records-file` and `--record` flags. These flags can be
// added through [addRecordsFlags].
// To get the records use [recordsFromFlags].
func recordsFromFlags(cmd *cobra.Command) ([]hcloud.ZoneRRSetRecord, error) {
	var parsedRecords []hcloud.ZoneRRSetRecord

	if cmd.Flags().Changed("records-file") {
		recordsFile, err := cmd.Flags().GetString("records-file")
		if err != nil {
			return nil, err
		}

		parsedRecords, err = parseRecords(recordsFile)
		if err != nil {
			return nil, err
		}
	} else {
		records, err := cmd.Flags().GetStringArray("record")
		if err != nil {
			return nil, err
		}

		parsedRecords = make([]hcloud.ZoneRRSetRecord, 0, len(records))
		for _, record := range records {
			parsedRecords = append(parsedRecords, hcloud.ZoneRRSetRecord{
				Value: record,
			})
		}
	}

	return parsedRecords, nil
}

// rrsetArgumentsCompletionFuncs provides completion funcs for the standard pattern:
// <zone> <name> <type>
// used by "hcloud zone rrset ..." commands.
func rrsetArgumentsCompletionFuncs(client hcapi2.Client) []cobra.CompletionFunc {

	return []cobra.CompletionFunc{
		cmpl.SuggestCandidatesF(client.Zone().Names),
		cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
			zoneIDOrName := args[0]
			zoneIDOrName, err := util.ParseZoneIDOrName(zoneIDOrName)
			if err != nil {
				return nil
			}

			rrsets, err := client.Zone().AllRRSets(context.Background(), &hcloud.Zone{Name: zoneIDOrName})
			if err != nil {
				return nil
			}

			uniqueRRSetNames := map[string]struct{}{}

			for _, rrset := range rrsets {
				uniqueRRSetNames[rrset.Name] = struct{}{}
			}

			return maps.Keys(uniqueRRSetNames)
		}),
		cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
			zoneIDOrName, rrsetName := args[0], args[1]
			zoneIDOrName, err := util.ParseZoneIDOrName(zoneIDOrName)
			if err != nil {
				return nil
			}

			rrsets, err := client.Zone().AllRRSetsWithOpts(context.Background(), &hcloud.Zone{Name: zoneIDOrName}, hcloud.ZoneRRSetListOpts{Name: rrsetName})
			if err != nil {
				return nil
			}

			uniqueRRSetTypes := map[string]struct{}{}

			for _, rrset := range rrsets {
				uniqueRRSetTypes[string(rrset.Type)] = struct{}{}
			}

			return maps.Keys(uniqueRRSetTypes)
		}),
	}
}

func FormatTXTRecords(cmd *cobra.Command, records []hcloud.ZoneRRSetRecord) {
	changed := false
	for i := range records {
		if !zoneutil.IsTXTRecordQuoted(records[i].Value) {
			records[i].Value = zoneutil.FormatTXTRecord(records[i].Value)
			cmd.Printf("Warning: Changed TXT record to: %s\n", records[i].Value)
			changed = true
		}
	}
	if changed {
		cmd.Printf("Warning: Learn more why in 'hcloud zone rrset --help'\n")
	}
}

package rrset

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

var recordsFileExample = `The optional records file has to be in JSON format. You can find the schema at https://docs.hetzner.cloud/reference/cloud#zone-rrset-actions-set-records-of-an-rrset

Example file content:

[
  {
    "value": "198.51.100.1",
    "comment": "My web server at Hetzner Cloud."
  },
  {
    "value": "198.51.100.2",
    "comment": "My other server at Hetzner Cloud."
  }
]`

var SetRecordsCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "set-records (--record <value>... | --records-file <file>) <zone> <name> <type>",
			Short: "Set the records of a Zone RRSet",
			Long: `Set the records of a Zone RRSet.

- If the Zone RRSet doesn't exist, it will be created.
- If the Zone RRSet already exists, its records will be replaced.
- If the provided records are empty, the Zone RRSet will be deleted.

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

		records, err := recordsFromFlags(cmd)
		if err != nil {
			return err
		}

		// TXT: Format record values to simplify its usage
		if hcloud.ZoneRRSetType(rrsetType) == hcloud.ZoneRRSetTypeTXT {
			FormatTXTRecords(cmd, records)
		}

		zoneIDOrName, err = util.ParseZoneIDOrName(zoneIDOrName)
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

		if len(records) == 0 {
			if rrset == nil {
				cmd.Printf("Zone RRSet %s %s doesn't exist. No action necessary.\n", rrsetName, rrsetType)
			} else {
				result, _, err := s.Client().Zone().DeleteRRSet(s, rrset)
				if err != nil {
					return err
				}
				if err := s.WaitForActions(s, cmd, result.Action); err != nil {
					return err
				}
				cmd.Printf("Zone RRSet %s %s deleted\n", rrset.Name, rrset.Type)
			}
			return nil
		}

		if rrset == nil {
			result, _, err := s.Client().Zone().CreateRRSet(s, zone, hcloud.ZoneRRSetCreateOpts{
				Name:    rrsetName,
				Type:    hcloud.ZoneRRSetType(rrsetType),
				Records: records,
			})
			if err != nil {
				return err
			}
			if err := s.WaitForActions(s, cmd, result.Action); err != nil {
				return err
			}
			rrset := result.RRSet
			cmd.Printf("Created and set records on Zone RRSet %s %s\n", rrset.Name, rrset.Type)
			return nil
		}

		opts := hcloud.ZoneRRSetSetRecordsOpts{Records: records}

		action, _, err := s.Client().Zone().SetRRSetRecords(s, rrset, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Set records on Zone RRSet %s %s\n", rrset.Name, rrset.Type)
		return nil
	},
	Experimental: experimental.DNS,
}

func parseRecords(path string) ([]hcloud.ZoneRRSetRecord, error) {
	var (
		data []byte
		err  error
	)
	if path == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return nil, err
	}

	var recordSchemas []schema.ZoneRRSetRecord
	err = json.Unmarshal(data, &recordSchemas)
	if err != nil {
		return nil, err
	}

	var records []hcloud.ZoneRRSetRecord
	for _, s := range recordSchemas {
		records = append(records, *hcloud.ZoneRRSetRecordFromSchema(s))
	}
	return records, err
}

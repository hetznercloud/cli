package rrset

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.ZoneRRSet]{
	ResourceNameSingular:       "Zone RRSet",
	ShortDescription:           "Describe a Zone RRSet",
	PositionalArgumentOverride: []string{"zone", "name", "type"},
	ValidArgsFunction:          rrsetArgumentsCompletionFuncs,
	FetchWithArgs: func(s state.State, _ *cobra.Command, args []string) (*hcloud.ZoneRRSet, interface{}, error) {
		zoneIDOrName, rrsetName, rrsetType := args[0], args[1], args[2]

		zoneIDOrName, err := util.ParseZoneIDOrName(zoneIDOrName)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert Zone name to ascii: %w", err)
		}

		zone, _, err := s.Client().Zone().Get(s, zoneIDOrName)
		if err != nil {
			return nil, nil, err
		}
		if zone == nil {
			return nil, nil, fmt.Errorf("Zone not found: %s", zoneIDOrName)
		}

		rrset, _, err := s.Client().Zone().GetRRSetByNameAndType(s, zone, rrsetName, hcloud.ZoneRRSetType(rrsetType))
		if err != nil {
			return nil, nil, err
		}
		return rrset, hcloud.SchemaFromZoneRRSet(rrset), nil
	},
	PrintText: func(_ state.State, cmd *cobra.Command, rrset *hcloud.ZoneRRSet) error {
		ttl := "-"
		if rrset.TTL != nil {
			ttl = strconv.Itoa(*rrset.TTL)
		}

		cmd.Printf("ID:\t\t%s\n", rrset.ID)
		cmd.Printf("Type:\t\t%s\n", rrset.Type)
		cmd.Printf("Name:\t\t%s\n", rrset.Name)
		cmd.Printf("TTL:\t\t%s\n", ttl)
		cmd.Printf("Protection:\n")
		cmd.Printf("  Change:\t%s\n", util.YesNo(rrset.Protection.Change))

		cmd.Print("Labels:\n")
		if len(rrset.Labels) == 0 {
			cmd.Print("  No labels\n")
		} else {
			for key, value := range util.IterateInOrder(rrset.Labels) {
				cmd.Printf("  %s: %s\n", key, value)
			}
		}

		cmd.Printf("Records:\n")
		if len(rrset.Records) == 0 {
			cmd.Print("  No Records\n")
		} else {
			for _, record := range rrset.Records {
				cmd.Printf("  - Value:\t%s\n", record.Value)
				if record.Comment != "" {
					cmd.Printf("    Comment:\t%s\n", record.Comment)
				}
			}
		}

		return nil
	},
	Experimental: experimental.DNS,
}

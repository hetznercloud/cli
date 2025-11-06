package rrset

import (
	"fmt"
	"io"
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
	PrintText: func(_ state.State, _ *cobra.Command, out io.Writer, rrset *hcloud.ZoneRRSet) error {
		ttl := "-"
		if rrset.TTL != nil {
			ttl = strconv.Itoa(*rrset.TTL)
		}

		fmt.Fprintf(out, "ID:\t%s\n", rrset.ID)
		fmt.Fprintf(out, "Type:\t%s\n", rrset.Type)
		fmt.Fprintf(out, "Name:\t%s\n", rrset.Name)
		fmt.Fprintf(out, "TTL:\t%s\n", ttl)

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Protection:\n")
		fmt.Fprintf(out, "  Change:\t%s\n", util.YesNo(rrset.Protection.Change))

		fmt.Fprintln(out)
		util.DescribeLabels(out, rrset.Labels, "")

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Records:\n")
		if len(rrset.Records) == 0 {
			fmt.Fprintf(out, "  No Records\n")
		} else {
			for _, record := range rrset.Records {
				fmt.Fprintf(out, "  - Value:\t%s\n", record.Value)
				if record.Comment != "" {
					fmt.Fprintf(out, "    Comment:\t%s\n", record.Comment)
				}
			}
		}

		return nil
	},
	Experimental: experimental.DNS,
}

package rrset

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var rrsetTypeStrings = []string{
	string(hcloud.ZoneRRSetTypeA), string(hcloud.ZoneRRSetTypeAAAA), string(hcloud.ZoneRRSetTypeCAA), string(hcloud.ZoneRRSetTypeCNAME),
	string(hcloud.ZoneRRSetTypeDS), string(hcloud.ZoneRRSetTypeHINFO), string(hcloud.ZoneRRSetTypeHTTPS), string(hcloud.ZoneRRSetTypeMX),
	string(hcloud.ZoneRRSetTypeNS), string(hcloud.ZoneRRSetTypePTR), string(hcloud.ZoneRRSetTypeRP), string(hcloud.ZoneRRSetTypeSOA),
	string(hcloud.ZoneRRSetTypeSRV), string(hcloud.ZoneRRSetTypeSVCB), string(hcloud.ZoneRRSetTypeTLSA), string(hcloud.ZoneRRSetTypeTXT),
}

var CreateCmd = base.CreateCmd[*hcloud.ZoneRRSet]{
	BaseCobraCommand: func(c hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   fmt.Sprintf("create [options] --name <name> --type <%s> (--record <record>... | --records-file <file>) <zone>", strings.Join(rrsetTypeStrings, "|")),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(c.Zone().Names)),
			Short:                 "Create a Zone RRSet",
			Long:                  "Create a Zone RRSet.\n\n" + recordsFileExample,
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("name", "", "Name of the Zone RRSet (required)")
		_ = cmd.MarkFlagRequired("name")

		cmd.Flags().String("type", "", "Type of the Zone RRSet (required)")
		_ = cmd.RegisterFlagCompletionFunc("type", cmpl.SuggestCandidates(rrsetTypeStrings...))
		_ = cmd.MarkFlagRequired("type")

		cmd.Flags().Int("ttl", 0, "Time To Live (TTL) of the RRSet")

		cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

		addRecordsFlags(cmd)

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) (*hcloud.ZoneRRSet, any, error) {
		zoneIDOrName := args[0]
		zoneIDOrName, err := util.ParseZoneIDOrName(zoneIDOrName)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert Zone name to ascii: %w", err)
		}

		name, _ := cmd.Flags().GetString("name")
		typ, _ := cmd.Flags().GetString("type")
		labels, _ := cmd.Flags().GetStringToString("label")

		if !slices.Contains(rrsetTypeStrings, typ) {
			return nil, nil, fmt.Errorf("unknown Zone RRSet type: %s", typ)
		}

		zone, _, err := s.Client().Zone().Get(s, zoneIDOrName)
		if err != nil {
			return nil, nil, err
		}
		if zone == nil {
			return nil, nil, fmt.Errorf("Zone not found: %s", zoneIDOrName)
		}

		createOpts := hcloud.ZoneRRSetCreateOpts{
			Name:   name,
			Type:   hcloud.ZoneRRSetType(typ),
			Labels: labels,
		}

		createOpts.Records, err = recordsFromFlags(cmd)
		if err != nil {
			return nil, nil, err
		}

		// TXT: Format record values to simplify its usage
		if createOpts.Type == hcloud.ZoneRRSetTypeTXT {
			FormatTXTRecords(cmd, createOpts.Records)
		}

		if cmd.Flags().Changed("ttl") {
			ttl, _ := cmd.Flags().GetInt("ttl")
			createOpts.TTL = &ttl
		}

		result, _, err := s.Client().Zone().CreateRRSet(s, zone, createOpts)
		if err != nil {
			return nil, nil, err
		}

		if err := s.WaitForActions(s, cmd, result.Action); err != nil {
			return nil, nil, err
		}
		cmd.Printf("Zone RRSet %s %s created\n", result.RRSet.Name, result.RRSet.Type)

		return result.RRSet, util.Wrap("rrset", hcloud.SchemaFromZoneRRSet(result.RRSet)), nil
	},
}

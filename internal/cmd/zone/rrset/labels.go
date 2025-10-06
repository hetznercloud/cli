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

var LabelCmds = base.LabelCmds[*hcloud.ZoneRRSet]{
	ResourceNameSingular:   "Zone RRSet",
	ShortDescriptionAdd:    "Add a label to a Zone RRSet",
	ShortDescriptionRemove: "Remove a label from a Zone RRSet",

	PositionalArgumentOverride: []string{"zone", "name", "type"},
	ValidArgsFunction: func(client hcapi2.Client) []cobra.CompletionFunc {
		validArgsFunction := rrsetArgumentsCompletionFuncs(client)

		validArgsFunction = append(validArgsFunction, cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
			if len(args) < 3 {
				return nil
			}
			zoneIDOrName, rrsetName, rrsetType := args[0], args[1], args[2]
			return client.Zone().RRSetLabelKeys(zoneIDOrName, rrsetName, hcloud.ZoneRRSetType(rrsetType))
		}))

		return validArgsFunction
	},
	FetchWithArgs: func(s state.State, args []string) (*hcloud.ZoneRRSet, error) {
		zoneIDOrName, rrsetName, rrsetType := args[0], args[1], args[2]

		zoneIDOrName, err := util.ParseZoneIDOrName(zoneIDOrName)
		if err != nil {
			return nil, fmt.Errorf("failed to convert Zone name to ascii: %w", err)
		}

		zone, _, err := s.Client().Zone().Get(s, zoneIDOrName)
		if err != nil {
			return nil, err
		}
		if zone == nil {
			return nil, fmt.Errorf("Zone not found: %s", zoneIDOrName)
		}

		rrset, _, err := s.Client().Zone().GetRRSetByNameAndType(s, zone, rrsetName, hcloud.ZoneRRSetType(rrsetType))
		if err != nil {
			return nil, err
		}
		if rrset == nil {
			return nil, fmt.Errorf("Zone RRSet not found: %s %s", rrsetName, rrsetType)
		}

		return rrset, nil
	},
	SetLabels: func(s state.State, rrset *hcloud.ZoneRRSet, labels map[string]string) error {
		opts := hcloud.ZoneRRSetUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Zone().UpdateRRSet(s, rrset, opts)
		return err
	},
	GetLabels: func(rrset *hcloud.ZoneRRSet) map[string]string {
		return rrset.Labels
	},
	GetIDOrName: func(rrset *hcloud.ZoneRRSet) string {
		return fmt.Sprintf("%s %s", rrset.Name, rrset.Type)
	},
	Experimental: experimental.DNS,
}

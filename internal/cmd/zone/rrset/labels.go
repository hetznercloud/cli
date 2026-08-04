package rrset

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
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

		validArgsFunction = append(validArgsFunction, cmpl.SuggestCandidatesCtxE(func(cmd *cobra.Command, args []string) ([]string, error) {
			if len(args) < 3 {
				return nil, nil
			}
			zoneIDOrName, rrsetName, rrsetType := args[0], args[1], args[2]
			return client.Zone().RRSetLabelKeys(cmd.Context(), zoneIDOrName, rrsetName, hcloud.ZoneRRSetType(rrsetType))
		}))

		return validArgsFunction
	},
	FetchWithArgs: func(ctx context.Context, s state.State, args []string) (*hcloud.ZoneRRSet, error) {
		zoneIDOrName, rrsetName, rrsetType := args[0], args[1], args[2]

		zoneIDOrName, err := util.ParseZoneIDOrName(zoneIDOrName)
		if err != nil {
			return nil, fmt.Errorf("failed to convert Zone name to ascii: %w", err)
		}

		zone, _, err := s.Client().Zone().Get(ctx, zoneIDOrName)
		if err != nil {
			return nil, err
		}
		if zone == nil {
			return nil, fmt.Errorf("Zone not found: %s", zoneIDOrName)
		}

		rrset, _, err := s.Client().Zone().GetRRSetByNameAndType(ctx, zone, rrsetName, hcloud.ZoneRRSetType(rrsetType))
		if err != nil {
			return nil, err
		}
		if rrset == nil {
			return nil, fmt.Errorf("Zone RRSet not found: %s %s", rrsetName, rrsetType)
		}

		return rrset, nil
	},
	SetLabels: func(ctx context.Context, s state.State, rrset *hcloud.ZoneRRSet, labels map[string]string) error {
		opts := hcloud.ZoneRRSetUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Zone().UpdateRRSet(ctx, rrset, opts)
		return err
	},
	GetLabels: func(rrset *hcloud.ZoneRRSet) map[string]string {
		return rrset.Labels
	},
	GetIDOrName: func(rrset *hcloud.ZoneRRSet) string {
		return fmt.Sprintf("%s %s", rrset.Name, rrset.Type)
	},
}

package rrset

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ChangeProtectionCmds = base.ChangeProtectionCmds[*hcloud.ZoneRRSet, hcloud.ZoneRRSetChangeProtectionOpts]{
	ResourceNameSingular: "Zone RRSet",

	ValidArgsFunction: rrsetArgumentsCompletionFuncs,

	PositionalArgumentOverride: []string{"zone", "name", "type"},

	ProtectionLevels: map[string]func(opts *hcloud.ZoneRRSetChangeProtectionOpts, value bool){
		"change": func(opts *hcloud.ZoneRRSetChangeProtectionOpts, value bool) {
			opts.Change = &value
		},
	},

	FetchWithArgs: func(s state.State, cmd *cobra.Command, args []string) (*hcloud.ZoneRRSet, *hcloud.Response, error) {
		zoneIDOrName, rrsetName, rrsetType := args[0], args[1], args[2]

		zoneIDOrName, err := util.ParseZoneIDOrName(zoneIDOrName)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert Zone name to ascii: %w", err)
		}

		zone, res, err := s.Client().Zone().Get(cmd.Context(), zoneIDOrName)
		if err != nil {
			return nil, res, err
		}
		if zone == nil {
			return nil, res, fmt.Errorf("Zone not found: %s", zoneIDOrName)
		}

		rrset, res, err := s.Client().Zone().GetRRSetByNameAndType(cmd.Context(), zone, rrsetName, hcloud.ZoneRRSetType(rrsetType))
		if err != nil {
			return nil, res, err
		}
		if rrset == nil {
			return nil, res, fmt.Errorf("Zone RRSet not found: %s %s", rrsetName, rrsetType)
		}

		return rrset, res, nil
	},

	ChangeProtectionFunction: func(ctx context.Context, s state.State, rrset *hcloud.ZoneRRSet, opts hcloud.ZoneRRSetChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error) {
		return s.Client().Zone().ChangeRRSetProtection(ctx, rrset, opts)
	},

	IDOrName: func(rrset *hcloud.ZoneRRSet) string {
		return fmt.Sprintf("%s %s", rrset.Name, rrset.Type)
	},
}

package rrset

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func getChangeProtectionOpts(enable bool, flags []string) (hcloud.ZoneRRSetChangeProtectionOpts, error) {
	opts := hcloud.ZoneRRSetChangeProtectionOpts{}

	var unknown []string
	for _, arg := range flags {
		switch strings.ToLower(arg) {
		case "change":
			opts.Change = &enable
		default:
			unknown = append(unknown, arg)
		}
	}
	if len(unknown) > 0 {
		return opts, fmt.Errorf("unknown protection level: %s", strings.Join(unknown, ", "))
	}

	return opts, nil
}

func changeProtection(s state.State, cmd *cobra.Command, rrset *hcloud.ZoneRRSet, enable bool, opts hcloud.ZoneRRSetChangeProtectionOpts) error {
	if opts.Change == nil {
		return nil
	}

	action, _, err := s.Client().Zone().ChangeRRSetProtection(s, rrset, opts)
	if err != nil {
		return err
	}

	if err = s.WaitForActions(s, cmd, action); err != nil {
		return err
	}

	if enable {
		cmd.Printf("Resource protection enabled for Zone RRSet %s %s\n", rrset.Name, rrset.Type)
	} else {
		cmd.Printf("Resource protection disabled for Zone RRSet %s %s\n", rrset.Name, rrset.Type)
	}
	return nil
}

var EnableProtectionCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:                   "enable-protection <zone> <name> <type> change",
			Args:                  util.ValidateLenient,
			Short:                 "Enable resource protection for a Zone RRSet",
			ValidArgsFunction:     cmpl.SuggestArgs(rrsetArgumentsCompletionFuncs(client)...),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
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

		rrset, _, err := s.Client().Zone().GetRRSetByNameAndType(s, zone, rrsetName, hcloud.ZoneRRSetType(rrsetType))
		if err != nil {
			return err
		}
		if rrset == nil {
			return fmt.Errorf("Zone RRSet not found: %s %s", rrsetName, rrsetType)
		}

		opts, err := getChangeProtectionOpts(true, args[3:])
		if err != nil {
			return err
		}

		return changeProtection(s, cmd, rrset, true, opts)
	},
}

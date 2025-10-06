package zone

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func getChangeProtectionOpts(enable bool, flags []string) (hcloud.ZoneChangeProtectionOpts, error) {
	opts := hcloud.ZoneChangeProtectionOpts{}

	var unknown []string
	for _, arg := range flags {
		switch strings.ToLower(arg) {
		case "delete":
			opts.Delete = &enable
		default:
			unknown = append(unknown, arg)
		}
	}
	if len(unknown) > 0 {
		return opts, fmt.Errorf("unknown protection level: %s", strings.Join(unknown, ", "))
	}

	return opts, nil
}

func changeProtection(s state.State, cmd *cobra.Command,
	zone *hcloud.Zone, enable bool, opts hcloud.ZoneChangeProtectionOpts) error {
	if opts.Delete == nil {
		return nil
	}

	action, _, err := s.Client().Zone().ChangeProtection(s, zone, opts)
	if err != nil {
		return err
	}

	if err := s.WaitForActions(s, cmd, action); err != nil {
		return err
	}

	if enable {
		cmd.Printf("Resource protection enabled for Zone %s\n", zone.Name)
	} else {
		cmd.Printf("Resource protection disabled for Zone %s\n", zone.Name)
	}
	return nil
}

var EnableProtectionCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:   "enable-protection <zone> delete",
			Args:  util.ValidateLenient,
			Short: "Enable resource protection for a Zone",
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.Zone().Names),
				cmpl.SuggestCandidates("delete"),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName, levels := args[0], args[1:]

		idOrName, err := util.ParseZoneIDOrName(idOrName)
		if err != nil {
			return fmt.Errorf("failed to convert Zone name to ascii: %w", err)
		}

		zone, _, err := s.Client().Zone().Get(s, idOrName)
		if err != nil {
			return err
		}
		if zone == nil {
			return fmt.Errorf("Zone not found: %s", idOrName)
		}

		opts, err := getChangeProtectionOpts(true, levels)
		if err != nil {
			return err
		}

		return changeProtection(s, cmd, zone, true, opts)
	},
	Experimental: experimental.DNS,
}

package zone

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ChangeTTLCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "change-ttl --ttl <ttl> <zone>",
			Short:                 "Changes the default Time To Live (TTL) of a Zone",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Zone().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().Int("ttl", 3600, "Default Time To Live (TTL) of the Zone (required)")
		_ = cmd.MarkFlagRequired("ttl")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
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

		ttl, _ := cmd.Flags().GetInt("ttl")
		opts := hcloud.ZoneChangeTTLOpts{TTL: ttl}

		action, _, err := s.Client().Zone().ChangeTTL(s, zone, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Changed default TTL on Zone %s\n", zone.Name)
		return nil
	},
}

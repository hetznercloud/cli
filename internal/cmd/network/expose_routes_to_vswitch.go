package network

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

var ExposeRoutesToVSwitchCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "expose-routes-to-vswitch [--disable] <network>",
			Short:                 "Expose routes to connected vSwitch",
			Long:                  "Enabling this will expose routes to the connected vSwitch. Set the --disable flag to remove the exposed routes.",
			Args:                  util.Validate,
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Network().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().Bool("disable", false, "Remove any exposed routes from the connected vSwitch")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		network, _, err := s.Client().Network().Get(s, idOrName)
		if err != nil {
			return err
		}
		if network == nil {
			return fmt.Errorf("network not found: %s", idOrName)
		}

		disable, _ := cmd.Flags().GetBool("disable")
		opts := hcloud.NetworkUpdateOpts{
			ExposeRoutesToVSwitch: hcloud.Ptr(!disable),
		}

		_, _, err = s.Client().Network().Update(s, network, opts)
		if err != nil {
			return err
		}

		if disable {
			cmd.Printf("Exposing routes to connected vSwitch of network %s disabled\n", network.Name)
		} else {
			cmd.Printf("Exposing routes to connected vSwitch of network %s enabled\n", network.Name)
		}

		return nil
	},
}

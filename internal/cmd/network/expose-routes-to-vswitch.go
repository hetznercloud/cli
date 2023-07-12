package network

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"
)

var ExposeRoutesToVSwitchCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "expose-routes-to-vswitch [flags] network",
			Short:                 "Expose routes to connected vSwitch",
			Long:                  "Enabling this will expose routes to the connected vSwitch. Set the --disable flag to remove the exposed routes.",
			Args:                  cobra.ExactArgs(1),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Network().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().Bool("disable", false, "Remove any exposed routes from the connected vSwitch")

		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		network, _, err := client.Network().Get(ctx, idOrName)
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

		_, _, err = client.Network().Update(ctx, network, opts)
		if err != nil {
			return err
		}

		if disable {
			fmt.Printf("Exposing routes to connected vSwitch of network %s disabled\n", network.Name)
		} else {
			fmt.Printf("Exposing routes to connected vSwitch of network %s enabled\n", network.Name)
		}

		return nil
	},
}

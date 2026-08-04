package server

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

var DetachFromNetworkCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "detach-from-network --network <network> <server>",
			Short:                 "Detach a Server from a Network",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.Server().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
		cmd.Flags().StringP("network", "n", "", "Network (ID or name) (required)")
		cmpl.RegisterFlagCompletion(cmd, "network", cmpl.SuggestCandidatesF(client.Network().Names))
		util.MarkFlagRequired(cmd, "network")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		server, _, err := s.Client().Server().Get(cmd.Context(), idOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("Server not found: %s", idOrName)
		}
		networkIDOrName, _ := cmd.Flags().GetString("network")
		network, _, err := s.Client().Network().Get(cmd.Context(), networkIDOrName)
		if err != nil {
			return err
		}
		if network == nil {
			return fmt.Errorf("Network not found: %s", networkIDOrName)
		}

		opts := hcloud.ServerDetachFromNetworkOpts{
			Network: network,
		}
		action, _, err := s.Client().Server().DetachFromNetwork(cmd.Context(), server, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(cmd.Context(), cmd, action); err != nil {
			return err
		}

		cmd.Printf("Server %d detached from Network %d\n", server.ID, network.ID)
		return nil
	},
}

package loadbalancer

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var AttachToNetworkCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "attach-to-network [--ip <ip>] --network <network> <load-balancer>",
			Short:                 "Attach a Load Balancer to a Network",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.LoadBalancer().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().StringP("network", "n", "", "Network (ID or name) (required)")
		_ = cmd.RegisterFlagCompletionFunc("network", cmpl.SuggestCandidatesF(client.Network().Names))
		_ = cmd.MarkFlagRequired("network")

		cmd.Flags().IP("ip", nil, "IP address to assign to the Load Balancer (auto-assigned if omitted)")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		loadBalancer, _, err := s.Client().LoadBalancer().Get(s, idOrName)
		if err != nil {
			return err
		}
		if loadBalancer == nil {
			return fmt.Errorf("Load Balancer not found: %s", idOrName)
		}

		networkIDOrName, _ := cmd.Flags().GetString("network")
		network, _, err := s.Client().Network().Get(s, networkIDOrName)
		if err != nil {
			return err
		}
		if network == nil {
			return fmt.Errorf("network not found: %s", networkIDOrName)
		}

		ip, _ := cmd.Flags().GetIP("ip")

		opts := hcloud.LoadBalancerAttachToNetworkOpts{
			Network: network,
			IP:      ip,
		}
		action, _, err := s.Client().LoadBalancer().AttachToNetwork(s, loadBalancer, opts)

		if err != nil {
			return err
		}

		if err := s.WaitForActions(s, cmd, action); err != nil {
			return err
		}

		cmd.Printf("Load Balancer %d attached to network %d\n", loadBalancer.ID, network.ID)
		return nil
	},
}

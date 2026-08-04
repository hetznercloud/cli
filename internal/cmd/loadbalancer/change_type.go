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

var ChangeTypeCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		return &cobra.Command{
			Use:   "change-type <load-balancer> <load-balancer-type>",
			Short: "Change type of a Load Balancer",
			ValidArgsFunction: cmpl.SuggestArgs(
				cmpl.SuggestCandidatesF(client.LoadBalancer().Names),
				cmpl.SuggestCandidatesF(client.LoadBalancerType().Names),
			),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		loadBalancer, _, err := s.Client().LoadBalancer().Get(cmd.Context(), idOrName)
		if err != nil {
			return err
		}
		if loadBalancer == nil {
			return fmt.Errorf("Load Balancer not found: %s", idOrName)
		}

		loadBalancerTypeIDOrName := args[1]
		loadBalancerType, _, err := s.Client().LoadBalancerType().Get(cmd.Context(), loadBalancerTypeIDOrName)
		if err != nil {
			return err
		}
		if loadBalancerType == nil {
			return fmt.Errorf("Load Balancer Type not found: %s", loadBalancerTypeIDOrName)
		}

		cmd.Print(deprecatedLoadBalancerTypeWarning(loadBalancerType))

		opts := hcloud.LoadBalancerChangeTypeOpts{
			LoadBalancerType: loadBalancerType,
		}
		action, _, err := s.Client().LoadBalancer().ChangeType(cmd.Context(), loadBalancer, opts)
		if err != nil {
			return err
		}

		if err := s.WaitForActions(cmd.Context(), cmd, action); err != nil {
			return err
		}

		cmd.Printf("Load Balancer %d changed to type %s\n", loadBalancer.ID, loadBalancerType.Name)
		return nil
	},
}

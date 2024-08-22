package loadbalancer

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var RemoveTargetCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "remove-target [options] <load-balancer>",
			Short:                 "Remove a target from a Load Balancer",
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.LoadBalancer().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().String("server", "", "Name or ID of the server")
		cmd.RegisterFlagCompletionFunc("server", cmpl.SuggestCandidatesF(client.Server().Names))

		cmd.Flags().String("label-selector", "", "Label Selector")

		cmd.Flags().String("ip", "", "IP address of an IP target")

		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		var (
			action       *hcloud.Action
			loadBalancer *hcloud.LoadBalancer
			err          error
		)

		serverIDOrName, _ := cmd.Flags().GetString("server")
		labelSelector, _ := cmd.Flags().GetString("label-selector")
		ipAddr, _ := cmd.Flags().GetString("ip")

		idOrName := args[0]

		loadBalancer, _, err = s.Client().LoadBalancer().Get(s, idOrName)
		if err != nil {
			return err
		}
		if loadBalancer == nil {
			return fmt.Errorf("Load Balancer not found: %s", idOrName)
		}

		if !util.AnySet(serverIDOrName, labelSelector, ipAddr) {
			return fmt.Errorf("specify one of --server, --label-selector, or --ip")
		}
		if !util.ExactlyOneSet(serverIDOrName, labelSelector, ipAddr) {
			return fmt.Errorf("--server, --label-selector, and --ip are mutually exclusive")
		}
		switch {
		case serverIDOrName != "":
			server, _, err := s.Client().Server().Get(s, serverIDOrName)
			if err != nil {
				return err
			}
			if server == nil {
				return fmt.Errorf("server not found: %s", serverIDOrName)
			}
			action, _, err = s.Client().LoadBalancer().RemoveServerTarget(s, loadBalancer, server)
			if err != nil {
				return err
			}
		case labelSelector != "":
			action, _, err = s.Client().LoadBalancer().RemoveLabelSelectorTarget(s, loadBalancer, labelSelector)
			if err != nil {
				return err
			}
		case ipAddr != "":
			ip := net.ParseIP(ipAddr)
			if ip == nil {
				return fmt.Errorf("invalid ip provided")
			}
			if action, _, err = s.Client().LoadBalancer().RemoveIPTarget(s, loadBalancer, ip); err != nil {
				return err
			}
		default:
			return fmt.Errorf("specify one of --server, --label-selector, or --ip")
		}

		if err := s.WaitForActions(cmd, s, action); err != nil {
			return err
		}
		cmd.Printf("Target removed from Load Balancer %d\n", loadBalancer.ID)

		return nil
	},
}

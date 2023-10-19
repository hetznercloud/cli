package loadbalancer

import (
	"context"
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

var RemoveTargetCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "remove-target LOADBALANCER FLAGS",
			Short:                 "Remove a target from a Load Balancer",
			Args:                  cobra.ExactArgs(1),
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
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		var (
			action       *hcloud.Action
			loadBalancer *hcloud.LoadBalancer
			err          error
		)

		serverIDOrName, _ := cmd.Flags().GetString("server")
		labelSelector, _ := cmd.Flags().GetString("label-selector")
		ipAddr, _ := cmd.Flags().GetString("ip")

		idOrName := args[0]

		loadBalancer, _, err = client.LoadBalancer().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if loadBalancer == nil {
			return fmt.Errorf("Load Balancer not found: %s", idOrName)
		}

		if !util.ExactlyOneSet(serverIDOrName, labelSelector, ipAddr) {
			return fmt.Errorf("--server, --label-selector, and --ip are mutually exclusive")
		}
		switch {
		case serverIDOrName != "":
			server, _, err := client.Server().Get(ctx, serverIDOrName)
			if err != nil {
				return err
			}
			if server == nil {
				return fmt.Errorf("server not found: %s", serverIDOrName)
			}
			action, _, err = client.LoadBalancer().RemoveServerTarget(ctx, loadBalancer, server)
			if err != nil {
				return err
			}
		case labelSelector != "":
			action, _, err = client.LoadBalancer().RemoveLabelSelectorTarget(ctx, loadBalancer, labelSelector)
			if err != nil {
				return err
			}
		case ipAddr != "":
			ip := net.ParseIP(ipAddr)
			if ip == nil {
				return fmt.Errorf("invalid ip provided")
			}
			if action, _, err = client.LoadBalancer().RemoveIPTarget(ctx, loadBalancer, ip); err != nil {
				return err
			}
		default:
			return fmt.Errorf("specify one of --server, --label-selector, or --ip")
		}

		if err := waiter.ActionProgress(ctx, action); err != nil {
			return err
		}
		fmt.Printf("Target removed from Load Balancer %d\n", loadBalancer.ID)

		return nil
	},
}

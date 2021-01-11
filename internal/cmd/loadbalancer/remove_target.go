package loadbalancer

import (
	"fmt"
	"net"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newRemoveTargetCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "remove-target LOADBALANCER FLAGS",
		Short:                 "Remove a target to a Load Balancer",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.LoadBalancerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runLoadBalancerRemoveTarget),
	}

	cmd.Flags().String("server", "", "Name or ID of the server")
	cmd.RegisterFlagCompletionFunc("server", cmpl.SuggestCandidatesF(cli.ServerNames))

	cmd.Flags().String("label-selector", "", "Label Selector")

	cmd.Flags().String("ip", "", "IP address of an IP target")

	return cmd
}

func runLoadBalancerRemoveTarget(cli *state.State, cmd *cobra.Command, args []string) error {
	var (
		action       *hcloud.Action
		loadBalancer *hcloud.LoadBalancer
		err          error
	)

	serverIDOrName, _ := cmd.Flags().GetString("server")
	labelSelector, _ := cmd.Flags().GetString("label-selector")
	ipAddr, _ := cmd.Flags().GetString("ip")

	idOrName := args[0]

	loadBalancer, _, err = cli.Client().LoadBalancer.Get(cli.Context, idOrName)
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
		server, _, err := cli.Client().Server.Get(cli.Context, serverIDOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", serverIDOrName)
		}
		action, _, err = cli.Client().LoadBalancer.RemoveServerTarget(cli.Context, loadBalancer, server)
		if err != nil {
			return err
		}
	case labelSelector != "":
		action, _, err = cli.Client().LoadBalancer.RemoveLabelSelectorTarget(cli.Context, loadBalancer, labelSelector)
		if err != nil {
			return err
		}
	case ipAddr != "":
		ip := net.ParseIP(ipAddr)
		if ip == nil {
			return fmt.Errorf("invalid ip provided")
		}
		if action, _, err = cli.Client().LoadBalancer.RemoveIPTarget(cli.Context, loadBalancer, ip); err != nil {
			return err
		}
	default:
		return fmt.Errorf("specify one of --server, --label-selector, or --ip")
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}
	fmt.Printf("Target removed from Load Balancer %d\n", loadBalancer.ID)

	return nil
}

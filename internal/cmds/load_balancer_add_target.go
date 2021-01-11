package cmds

import (
	"fmt"
	"net"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newLoadBalancerAddTargetCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add-target LOADBALANCER FLAGS",
		Short:                 "Add a target to a Load Balancer",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.LoadBalancerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runLoadBalancerAddTarget),
	}

	cmd.Flags().String("server", "", "Name or ID of the server")
	cmd.RegisterFlagCompletionFunc("server", cmpl.SuggestCandidatesF(cli.ServerNames))

	cmd.Flags().String("label-selector", "", "Label Selector")

	cmd.Flags().Bool("use-private-ip", false, "Determine if the Load Balancer should connect to the target via the network")
	cmd.Flags().String("ip", "", "Use the passed IP address as target")
	return cmd
}

func runLoadBalancerAddTarget(cli *state.State, cmd *cobra.Command, args []string) error {
	var (
		action       *hcloud.Action
		loadBalancer *hcloud.LoadBalancer
		err          error
	)

	idOrName := args[0]
	usePrivateIP, _ := cmd.Flags().GetBool("use-private-ip")
	serverIDOrName, _ := cmd.Flags().GetString("server")
	labelSelector, _ := cmd.Flags().GetString("label-selector")
	ipAddr, _ := cmd.Flags().GetString("ip")

	if !util.ExactlyOneSet(serverIDOrName, labelSelector, ipAddr) {
		return fmt.Errorf("--server, --label-selector, and --ip are mutually exclusive")
	}
	if loadBalancer, _, err = cli.Client().LoadBalancer.Get(cli.Context, idOrName); err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("Load Balancer not found: %s", idOrName)
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
		action, _, err = cli.Client().LoadBalancer.AddServerTarget(cli.Context, loadBalancer, hcloud.LoadBalancerAddServerTargetOpts{
			Server:       server,
			UsePrivateIP: hcloud.Bool(usePrivateIP),
		})
		if err != nil {
			return err
		}
	case labelSelector != "":
		action, _, err = cli.Client().LoadBalancer.AddLabelSelectorTarget(cli.Context, loadBalancer, hcloud.LoadBalancerAddLabelSelectorTargetOpts{
			Selector:     labelSelector,
			UsePrivateIP: hcloud.Bool(usePrivateIP),
		})
		if err != nil {
			return err
		}
	case ipAddr != "":
		ip := net.ParseIP(ipAddr)
		if ip == nil {
			return fmt.Errorf("invalid ip provided")
		}
		action, _, err = cli.Client().LoadBalancer.AddIPTarget(cli.Context, loadBalancer, hcloud.LoadBalancerAddIPTargetOpts{
			IP: ip,
		})
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("specify one of --server, --label-selector, or --ip")
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}
	fmt.Printf("Target added to Load Balancer %d\n", loadBalancer.ID)

	return nil
}

package cli

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newLoadBalancerAddTargetCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add-target LOADBALANCER FLAGS",
		Short:                 "Add a target to a Load Balancer",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runLoadBalancerAddTarget),
	}

	cmd.Flags().String("server", "", "Name or ID of the server")
	cmd.Flag("server").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_server_names"},
	}
	cmd.Flags().Bool("use-private-ip", false, "Determine if the Load Balancer should connect to the target via the network")
	return cmd
}

func runLoadBalancerAddTarget(cli *CLI, cmd *cobra.Command, args []string) error {
	serverIdOrName, _ := cmd.Flags().GetString("server")
	idOrName := args[0]
	usePrivateIP, _ := cmd.Flags().GetBool("use-private-ip")

	loadBalancer, _, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("Load Balancer not found: %s", idOrName)
	}

	var action *hcloud.Action
	if serverIdOrName != "" {
		server, _, err := cli.Client().Server.Get(cli.Context, serverIdOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", serverIdOrName)
		}
		action, _, err = cli.Client().LoadBalancer.AddServerTarget(cli.Context, loadBalancer, hcloud.LoadBalancerAddServerTargetOpts{
			Server:       server,
			UsePrivateIP: hcloud.Bool(usePrivateIP),
		})
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("specify one of server")
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}
	fmt.Printf("Target added to Load Balancer %d\n", loadBalancer.ID)

	return nil
}

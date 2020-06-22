package cli

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newLoadBalancerRemoveTargetCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "remove-target LOADBALANCER FLAGS",
		Short:                 "Remove a target to a Load Balancer",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runLoadBalancerRemoveTarget),
	}

	cmd.Flags().String("server", "", "Name or ID of the server")
	cmd.Flag("server").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_server_names"},
	}

	return cmd
}

func runLoadBalancerRemoveTarget(cli *CLI, cmd *cobra.Command, args []string) error {
	serverIdOrName, _ := cmd.Flags().GetString("server")
	idOrName := args[0]

	loadBalancer, _, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("Load Balancer not found: %s", idOrName)
	}

	var action *hcloud.Action
	if serverIdOrName == "" {
		return fmt.Errorf("specify a server")
	} else if serverIdOrName != "" {
		server, _, err := cli.Client().Server.Get(cli.Context, serverIdOrName)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %s", serverIdOrName)
		}
		action, _, err = cli.Client().LoadBalancer.RemoveServerTarget(cli.Context, loadBalancer, server)
		if err != nil {
			return err
		}
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}
	fmt.Printf("Target removed from Load Balancer %d\n", loadBalancer.ID)

	return nil
}

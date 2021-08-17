package loadbalancer

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newSetRDNSCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "set-rdns [FLAGS] LOADBALANCER",
		Short:                 "Change reverse DNS of a LOADBALANCER",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.LoadBalancerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runSetRDNS),
	}

	cmd.Flags().StringP("hostname", "r", "", "Hostname to set as a reverse DNS PTR entry (required)")
	cmd.MarkFlagRequired("hostname")

	cmd.Flags().StringP("ip", "i", "", "IP address for which the reverse DNS entry should be set")
	cmd.MarkFlagRequired("ip")
	return cmd
}

func runSetRDNS(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	loadBalancer, _, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("LoadBalancer not found: %v", idOrName)
	}

	ip, _ := cmd.Flags().GetString("ip")

	hostname, _ := cmd.Flags().GetString("hostname")
	action, _, err := cli.Client().LoadBalancer.ChangeDNSPtr(cli.Context, loadBalancer, ip, hcloud.String(hostname))
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Reverse DNS of Load Balancer %d changed\n", loadBalancer.ID)

	return nil
}

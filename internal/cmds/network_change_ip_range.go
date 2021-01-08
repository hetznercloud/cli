package cmds

import (
	"fmt"
	"net"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newNetworkChangeIPRangeCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "change-ip-range [FLAGS] NETWORK",
		Short:                 "Change the IP range of a network",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.NetworkNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runNetworkChangeIPRange),
	}

	cmd.Flags().IPNet("ip-range", net.IPNet{}, "New IP range (required)")
	cmd.MarkFlagRequired("ip-range")

	return cmd
}

func runNetworkChangeIPRange(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	network, _, err := cli.Client().Network.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if network == nil {
		return fmt.Errorf("network not found: %s", idOrName)
	}

	ipRange, _ := cmd.Flags().GetIPNet("ip-range")
	opts := hcloud.NetworkChangeIPRangeOpts{
		IPRange: &ipRange,
	}

	action, _, err := cli.Client().Network.ChangeIPRange(cli.Context, network, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}
	fmt.Printf("IP range of network %d changed\n", network.ID)
	return nil
}

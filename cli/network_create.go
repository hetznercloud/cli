package cli

import (
	"fmt"
	"net"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newNetworkCreateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create [FLAGS]",
		Short:                 "Create a network",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runNetworkCreate),
	}

	cmd.Flags().String("name", "", "Network name")
	cmd.MarkFlagRequired("name")

	cmd.Flags().IPNet("ip-range", net.IPNet{}, "Network IP range")
	cmd.MarkFlagRequired("ip-range")

	cmd.Flags().StringToString("label", nil, "User-defined labels ('key=value') (can be specified multiple times)")

	return cmd
}

func runNetworkCreate(cli *CLI, cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	ipRange, _ := cmd.Flags().GetIPNet("ip-range")
	labels, _ := cmd.Flags().GetStringToString("label")

	opts := hcloud.NetworkCreateOpts{
		Name:    name,
		IPRange: &ipRange,
		Labels:  labels,
	}

	network, _, err := cli.Client().Network.Create(cli.Context, opts)
	if err != nil {
		return err
	}

	fmt.Printf("Network %d created\n", network.ID)
	return nil
}

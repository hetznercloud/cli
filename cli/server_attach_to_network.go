package cli

import (
	"fmt"
	"net"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerAttachToNetworkCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "attach-to-network [FLAGS] SERVER",
		Short:                 "Attach a server to a network",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runServerAttachToNetwork),
	}

	cmd.Flags().StringP("network", "n", "", "Network (ID or name) (required)")
	cmd.Flag("network").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__hcloud_network_names"},
	}
	cmd.MarkFlagRequired("network")

	cmd.Flags().IP("ip", nil, "IP address to assign to the server (auto-assigned if omitted)")
	cmd.Flags().IPSlice("alias-ips", []net.IP{}, "Additional IP addresses to be assigned to the server")

	return cmd
}

func runServerAttachToNetwork(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	networkIDOrName, _ := cmd.Flags().GetString("network")
	network, _, err := cli.Client().Network.Get(cli.Context, networkIDOrName)
	if err != nil {
		return err
	}
	if network == nil {
		return fmt.Errorf("network not found: %s", networkIDOrName)
	}

	ip, _ := cmd.Flags().GetIP("ip")
	aliasIPs, _ := cmd.Flags().GetIPSlice("alias-ips")

	opts := hcloud.ServerAttachToNetworkOpts{
		Network: network,
		IP:      ip,
	}
	for _, aliasIP := range aliasIPs {
		opts.AliasIPs = append(opts.AliasIPs, aliasIP)
	}
	action, _, err := cli.Client().Server.AttachToNetwork(cli.Context, server, opts)

	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Server %d attached to network %d\n", server.ID, network.ID)
	return nil
}

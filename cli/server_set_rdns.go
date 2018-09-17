package cli

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerSetRDNSCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "set-rdns [FLAGS] SERVER",
		Short:                 "Change reverse DNS of a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runServerSetRDNS),
	}

	cmd.Flags().StringP("hostname", "r", "", "Hostname to set as a reverse DNS PTR entry")
	cmd.MarkFlagRequired("hostname")

	cmd.Flags().StringP("ip", "i", "", "IP address for which the reverse DNS entry should be set")

	return cmd
}

func runServerSetRDNS(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}
	ip, _ := cmd.Flags().GetString("ip")
	if ip == "" {
		ip = server.PublicNet.IPv4.IP.String()
	}

	hostname, _ := cmd.Flags().GetString("hostname")
	action, _, err := cli.Client().Server.ChangeDNSPtr(cli.Context, server, ip, hcloud.String(hostname))
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Reverse DNS of server %d changed\n", server.ID)

	return nil
}

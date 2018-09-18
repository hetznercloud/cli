package cli

import (
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newFloatingIPSetRDNSCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "set-rdns [FLAGS] FLOATINGIP",
		Short:                 "Change reverse DNS of a Floating IP",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runFloatingIPSetRDNS),
	}

	cmd.Flags().StringP("hostname", "r", "", "Hostname to set as a reverse DNS PTR entry")
	cmd.MarkFlagRequired("hostname")

	cmd.Flags().StringP("ip", "i", "", "IP address for which the reverse DNS entry should be set")
	return cmd
}

func runFloatingIPSetRDNS(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	floatingIP, _, err := cli.Client().FloatingIP.GetByID(cli.Context, id)
	if err != nil {
		return err
	}
	if floatingIP == nil {
		return fmt.Errorf("Floating IP not found: %d", id)
	}

	ip, _ := cmd.Flags().GetString("ip")
	if ip == "" {
		ip = floatingIP.IP.String()
	}

	hostname, _ := cmd.Flags().GetString("hostname")
	action, _, err := cli.Client().FloatingIP.ChangeDNSPtr(cli.Context, floatingIP, ip, hcloud.String(hostname))
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Reverse DNS of Floating IP %d changed\n", floatingIP.ID)

	return nil
}

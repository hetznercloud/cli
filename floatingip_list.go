package cli

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newFloatingIPListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "list",
		Short:            "List Floating IPs",
		TraverseChildren: true,
		RunE:             cli.wrap(runFloatingIPList),
	}
	return cmd
}

func runFloatingIPList(cli *CLI, cmd *cobra.Command, args []string) error {
	floatingIPs, err := cli.Client().FloatingIP.All(cli.Context)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tTYPE\tDESCRIPTION\tIP\tHOME\tSERVER\tDNS")
	for _, floatingIP := range floatingIPs {
		server := "-"
		if floatingIP.Server != nil {
			server = strconv.Itoa(floatingIP.Server.ID)
		}

		dns := "-"
		switch {
		case floatingIP.DNSPtr == nil, len(floatingIP.DNSPtr) == 0:
			dns = "-"
		case len(floatingIP.DNSPtr) == 1:
			for _, v := range floatingIP.DNSPtr {
				dns = v
				break
			}
		default:
			dns = fmt.Sprintf("%d entries", len(floatingIP.DNSPtr))
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n", floatingIP.ID, floatingIP.Type,
			floatingIP.Description, floatingIP.IP, floatingIP.HomeLocation.Name,
			server, dns)
	}
	w.Flush()

	return nil
}

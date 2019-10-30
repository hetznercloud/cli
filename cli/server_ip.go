package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newServerIPCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "ip SERVER FLAGS",
		Short:                 "Get the IP from a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runServerIP),
	}
	cmd.Flags().BoolP("ipv6", "6", false, "Return the first address of the IPv6 public server network")
	return cmd
}

func runServerIP(cli *CLI, cmd *cobra.Command, args []string) error {
	ipv6, err := cmd.Flags().GetBool("ipv6")
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}
	if ipv6 {
		fmt.Println(server.PublicNet.IPv6.IP.String() + "1")
	} else {
		fmt.Println(server.PublicNet.IPv4.IP.String())
	}
	return nil
}

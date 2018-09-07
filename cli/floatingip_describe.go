package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newFloatingIPDescribeCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] FLOATINGIP",
		Short:                 "Describe a Floating IP",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runFloatingIPDescribe),
	}
	return cmd
}

func runFloatingIPDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid Floating IP ID")
	}

	floatingIP, _, err := cli.Client().FloatingIP.GetByID(cli.Context, id)
	if err != nil {
		return err
	}
	if floatingIP == nil {
		return fmt.Errorf("Floating IP not found: %d", id)
	}

	fmt.Printf("ID:\t\t%d\n", floatingIP.ID)
	fmt.Printf("Type:\t\t%s\n", floatingIP.Type)
	fmt.Printf("Description:\t%s\n", na(floatingIP.Description))
	if floatingIP.Network != nil {
		fmt.Printf("IP:\t\t%s\n", floatingIP.Network.String())
	} else {
		fmt.Printf("IP:\t\t%s\n", floatingIP.IP.String())
	}
	fmt.Printf("Blocked:\t%s\n", yesno(floatingIP.Blocked))
	fmt.Printf("Home Location:\t%s\n", floatingIP.HomeLocation.Name)
	if floatingIP.Server != nil {
		server, _, err := cli.Client().Server.GetByID(cli.Context, floatingIP.Server.ID)
		if err != nil {
			return err
		}
		if server == nil {
			return fmt.Errorf("server not found: %d", id)
		}
		fmt.Printf("Server:\n")
		fmt.Printf("  ID:\t%d\n", server.ID)
		fmt.Printf("  Name:\t%s\n", server.Name)
	} else {
		fmt.Print("Server:\n  Not assigned\n")
	}
	fmt.Print("DNS:\n")
	if len(floatingIP.DNSPtr) == 0 {
		fmt.Print("  No reverse DNS entries\n")
	} else {
		for ip, dns := range floatingIP.DNSPtr {
			fmt.Printf("  %s: %s\n", ip, dns)
		}
	}

	fmt.Printf("Protection:\n")
	fmt.Printf("  Delete:\t%s\n", yesno(floatingIP.Protection.Delete))

	fmt.Print("Labels:\n")
	if len(floatingIP.Labels) == 0 {
		fmt.Print("  No labels\n")
	} else {
		for key, value := range floatingIP.Labels {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}
	return nil
}

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
		RunE: cli.wrap(runFloatingIPDescribe),
	}
	return cmd
}

func runFloatingIPDescribe(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid Floating IP ID")
	}

	floatingIP, _, err := cli.Client().FloatingIP.Get(cli.Context, id)
	if err != nil {
		return err
	}
	if floatingIP == nil {
		return fmt.Errorf("Floating IP not found: %d", id)
	}

	fmt.Printf("ID:\t\t%d\n", floatingIP.ID)
	fmt.Printf("Type:\t\t%s\n", floatingIP.Type)
	fmt.Printf("Description:\t%s\n", floatingIP.Description)
	fmt.Printf("Home Location:\t%s\n", floatingIP.HomeLocation.Name)
	if floatingIP.Server != nil {
		fmt.Printf("Server:\t\t%d\n", floatingIP.Server.ID)
	} else {
		fmt.Print("Server:\t\tnot assigned\n")
	}
	fmt.Print("DNS:\n")
	if len(floatingIP.DNSPtr) == 0 {
		fmt.Print("none set\n")
	} else {
		for ip, dns := range floatingIP.DNSPtr {
			fmt.Printf("  %s: %s\n", ip, dns)
		}
	}

	return nil
}

package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newFloatingIPAssignCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "assign [FLAGS] FLOATINGIP SERVER",
		Short:                 "Assign a Floating IP to a server",
		Args:                  cobra.ExactArgs(2),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runFloatingIPAssign),
	}
	cmd.MarkFlagRequired("server")
	return cmd
}

func runFloatingIPAssign(cli *CLI, cmd *cobra.Command, args []string) error {
	floatingIPID, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid Floating IP ID")
	}
	floatingIP := &hcloud.FloatingIP{ID: floatingIPID}

	serverIDOrName := args[1]
	server, _, err := cli.Client().Server.Get(cli.Context, serverIDOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", serverIDOrName)
	}

	action, _, err := cli.Client().FloatingIP.Assign(cli.Context, floatingIP, server)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Floating IP %d assigned to server %d\n", floatingIP.ID, server.ID)
	return nil
}

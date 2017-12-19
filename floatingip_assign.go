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
		RunE: cli.wrap(runFloatingIPAssign),
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

	serverID, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.New("invalid server ID")
	}
	server := &hcloud.Server{ID: serverID}

	action, _, err := cli.Client().FloatingIP.Assign(cli.Context, floatingIP, server)
	if err != nil {
		return err
	}
	errCh, _ := waitAction(cli.Context, cli.Client(), action)
	if err := <-errCh; err != nil {
		return err
	}
	fmt.Printf("Floating IP %d assigned to server %d\n", floatingIP.ID, server.ID)
	return nil
}

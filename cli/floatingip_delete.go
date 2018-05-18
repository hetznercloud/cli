package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newFloatingIPDeleteCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete [FLAGS] FLOATINGIP",
		Short:                 "Delete a Floating IP",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runFloatingIPDelete),
	}
	return cmd
}

func runFloatingIPDelete(cli *CLI, cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid Floating IP ID")
	}
	floatingIP := &hcloud.FloatingIP{ID: id}

	if _, err := cli.Client().FloatingIP.Delete(cli.Context, floatingIP); err != nil {
		return err
	}
	fmt.Printf("Floating IP %d deleted\n", id)
	return nil
}

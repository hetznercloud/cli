package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newFloatingIPUpdateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "update [FLAGS] FLOATINGIP",
		Short:                 "Update a Floating IP",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureActiveContext,
		RunE:                  cli.wrap(runFloatingIPUpdate),
	}

	cmd.Flags().String("description", "", "Floating IP description")

	return cmd
}

func runFloatingIPUpdate(cli *CLI, cmd *cobra.Command, args []string) error {
	floatingIPID, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid Floating IP ID")
	}
	floatingIP := &hcloud.FloatingIP{ID: floatingIPID}

	description, _ := cmd.Flags().GetString("description")
	opts := hcloud.FloatingIPUpdateOpts{
		Description: description,
	}
	_, _, err = cli.Client().FloatingIP.Update(cli.Context, floatingIP, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Floating IP %d updated\n", floatingIP.ID)
	return nil
}

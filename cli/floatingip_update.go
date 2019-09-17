package cli

import (
	"fmt"
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
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runFloatingIPUpdate),
	}

	cmd.Flags().String("description", "", "Floating IP description")
	cmd.Flags().String("name", "", "Floating IP name")

	return cmd
}

func runFloatingIPUpdate(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	floatingIP, _, err := cli.Client().FloatingIP.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if floatingIP == nil {
		return fmt.Errorf("Floating IP not found: %s", idOrName)
	}

	description, _ := cmd.Flags().GetString("description")
	name, _ := cmd.Flags().GetString("name")
	opts := hcloud.FloatingIPUpdateOpts{
		Description: description,
		Name:        name,
	}
	_, _, err = cli.Client().FloatingIP.Update(cli.Context, floatingIP, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Floating IP %d updated\n", floatingIP.ID)
	return nil
}

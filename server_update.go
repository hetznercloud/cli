package cli

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerUpdateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "update [FLAGS] SERVER",
		Short:                 "Update a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runServerUpdate),
	}

	cmd.Flags().String("name", "", "Server name")

	return cmd
}

func runServerUpdate(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	name, _ := cmd.Flags().GetString("name")
	opts := hcloud.ServerUpdateOpts{
		Name: name,
	}
	_, _, err = cli.Client().Server.Update(cli.Context, server, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Server %s updated\n", idOrName)
	return nil
}

package server

import (
	"errors"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newUpdateCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "update [FLAGS] SERVER",
		Short:                 "Update a server",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ServerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runUpdate),
	}

	cmd.Flags().String("name", "", "Server name")

	return cmd
}

func runUpdate(cli *state.State, cmd *cobra.Command, args []string) error {
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
	if opts.Name == "" {
		return errors.New("no updates")
	}

	_, _, err = cli.Client().Server.Update(cli.Context, server, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Server %d updated\n", server.ID)
	return nil
}

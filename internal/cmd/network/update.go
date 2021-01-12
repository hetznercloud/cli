package network

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
		Use:                   "update [FLAGS] NETWORK",
		Short:                 "Update a network",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.NetworkNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runUpdate),
	}

	cmd.Flags().String("name", "", "Network name")

	return cmd
}

func runUpdate(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	network, _, err := cli.Client().Network.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if network == nil {
		return fmt.Errorf("network not found: %s", idOrName)
	}

	name, _ := cmd.Flags().GetString("name")
	opts := hcloud.NetworkUpdateOpts{
		Name: name,
	}
	if opts.Name == "" {
		return errors.New("no updates")
	}

	_, _, err = cli.Client().Network.Update(cli.Context, network, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Network %d updated\n", network.ID)
	return nil
}

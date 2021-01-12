package server

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func newAttachISOCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "attach-iso [FLAGS] SERVER ISO",
		Short:            "Attach an ISO to a server",
		Args:             cobra.ExactArgs(2),
		TraverseChildren: true,
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.ServerNames),
			cmpl.SuggestCandidatesF(cli.ISONames),
		),
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runAttachISO),
	}

	return cmd
}

func runAttachISO(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	isoIDOrName := args[1]
	iso, _, err := cli.Client().ISO.Get(cli.Context, isoIDOrName)
	if err != nil {
		return err
	}
	if iso == nil {
		return fmt.Errorf("ISO not found: %s", isoIDOrName)
	}

	action, _, err := cli.Client().Server.AttachISO(cli.Context, server, iso)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("ISO %s attached to server %d\n", iso.Name, server.ID)
	return nil
}

package cmds

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerChangeTypeCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-type [FLAGS] SERVER SERVERTYPE",
		Short: "Change type of a server",
		Args:  cobra.ExactArgs(2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.ServerNames),
			cmpl.SuggestCandidatesF(cli.ServerTypeNames),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runServerChangeType),
	}

	cmd.Flags().Bool("keep-disk", false, "Keep disk size of current server type. This enables downgrading the server.")
	return cmd
}

func runServerChangeType(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	serverTypeIDOrName := args[1]
	serverType, _, err := cli.Client().ServerType.Get(cli.Context, serverTypeIDOrName)
	if err != nil {
		return err
	}
	if serverType == nil {
		return fmt.Errorf("server type not found: %s", serverTypeIDOrName)
	}

	keepDisk, _ := cmd.Flags().GetBool("keep-disk")
	opts := hcloud.ServerChangeTypeOpts{
		ServerType:  serverType,
		UpgradeDisk: !keepDisk,
	}
	action, _, err := cli.Client().Server.ChangeType(cli.Context, server, opts)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	if opts.UpgradeDisk {
		fmt.Printf("Server %d changed to type %s\n", server.ID, serverType.Name)
	} else {
		fmt.Printf("Server %d changed to type %s (disk size was unchanged)\n", server.ID, serverType.Name)
	}
	return nil
}

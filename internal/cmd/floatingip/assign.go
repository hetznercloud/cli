package floatingip

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func newAssignCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assign [FLAGS] FLOATINGIP SERVER",
		Short: "Assign a Floating IP to a server",
		Args:  cobra.ExactArgs(2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.FloatingIPNames),
			cmpl.SuggestCandidatesF(cli.ServerNames),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runAssign),
	}
	return cmd
}

func runAssign(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	floatingIP, _, err := cli.Client().FloatingIP.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if floatingIP == nil {
		return fmt.Errorf("Floating IP not found: %v", idOrName)
	}

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

package floatingip

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func newUnassignCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "unassign [FLAGS] FLOATINGIP",
		Short:                 "Unassign a Floating IP",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.FloatingIPNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runFloatingIPUnassign),
	}
	return cmd
}

func runFloatingIPUnassign(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	floatingIP, _, err := cli.Client().FloatingIP.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if floatingIP == nil {
		return fmt.Errorf("Floating IP not found: %v", idOrName)
	}

	action, _, err := cli.Client().FloatingIP.Unassign(cli.Context, floatingIP)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}

	fmt.Printf("Floating IP %d unassigned\n", floatingIP.ID)
	return nil
}

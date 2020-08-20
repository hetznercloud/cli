package cli

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/spf13/cobra"
)

func newFloatingIPDeleteCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete [FLAGS] FLOATINGIP",
		Short:                 "Delete a Floating IP",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.FloatingIPNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runFloatingIPDelete),
	}
	return cmd
}

func runFloatingIPDelete(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	floatingIP, _, err := cli.Client().FloatingIP.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if floatingIP == nil {
		return fmt.Errorf("Floating IP not found: %v", idOrName)
	}

	if _, err := cli.Client().FloatingIP.Delete(cli.Context, floatingIP); err != nil {
		return err
	}
	fmt.Printf("Floating IP %v deleted\n", idOrName)
	return nil
}

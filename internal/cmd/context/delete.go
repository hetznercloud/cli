package context

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
)

func newDeleteCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete [FLAGS] NAME",
		Short:                 "Delete a context",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.Config.ContextNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.Wrap(runDelete),
	}
	return cmd
}

func runDelete(cli *state.State, cmd *cobra.Command, args []string) error {
	name := args[0]
	context := cli.Config.ContextByName(name)
	if context == nil {
		return fmt.Errorf("context not found: %v", name)
	}
	if cli.Config.ActiveContext == context {
		_, _ = fmt.Fprintln(os.Stderr, "Warning: You are deleting the currently active context. Please select a new active context.")
		cli.Config.ActiveContext = nil
	}
	cli.Config.RemoveContext(context)
	return cli.WriteConfig()
}

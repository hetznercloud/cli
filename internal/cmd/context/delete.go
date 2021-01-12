package context

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
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
		cli.Config.ActiveContext = nil
	}
	cli.Config.RemoveContext(context)
	return cli.WriteConfig()
}

package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
)

func newUseCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "use [FLAGS] NAME",
		Short:                 "Use a context",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.Config.ContextNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.Wrap(runUse),
	}
	return cmd
}

func runUse(cli *state.State, cmd *cobra.Command, args []string) error {
	name := args[0]
	context := cli.Config.ContextByName(name)
	if context == nil {
		return fmt.Errorf("context not found: %v", name)
	}
	cli.Config.ActiveContext = context
	return cli.WriteConfig()
}

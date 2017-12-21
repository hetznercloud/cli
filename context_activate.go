package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newContextActivateCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "activate [FLAGS] NAME",
		Short:                 "Activate a context",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runContextActivate),
	}
	return cmd
}

func runContextActivate(cli *CLI, cmd *cobra.Command, args []string) error {
	name := args[0]
	context := cli.Config.ContextByName(name)
	if context == nil {
		return fmt.Errorf("context not found: %v", name)
	}
	cli.Config.ActiveContext = context
	return cli.WriteConfig()
}

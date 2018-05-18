package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newContextDeleteCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "delete [FLAGS] NAME",
		Short:                 "Delete a context",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runContextDelete),
	}
	return cmd
}

func runContextDelete(cli *CLI, cmd *cobra.Command, args []string) error {
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

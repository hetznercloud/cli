package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newContextActiveCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "active [FLAGS]",
		Short:                 "Show active context",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.wrap(runContextActive),
	}
	return cmd
}

func runContextActive(cli *CLI, cmd *cobra.Command, args []string) error {
	if cli.Config.ActiveContext != nil {
		fmt.Println(cli.Config.ActiveContext.Name)
	}
	return nil
}

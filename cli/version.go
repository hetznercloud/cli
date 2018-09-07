package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is set via compiler flags (see script/build.bash)
var Version = "was not built properly"

func newVersionCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "version",
		Short:                 "Print version information",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		RunE:                  cli.wrap(runVersion),
	}
	return cmd
}

func runVersion(cli *CLI, cmd *cobra.Command, args []string) error {
	fmt.Printf("hcloud %s\n", Version)
	return nil
}

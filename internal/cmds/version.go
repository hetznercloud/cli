package cmds

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

// Version is set via compiler flags (see script/build.bash)
var Version = "was not built properly"

func NewVersionCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "version",
		Short:                 "Print version information",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		RunE:                  cli.Wrap(runVersion),
	}
	return cmd
}

func runVersion(cli *state.State, cmd *cobra.Command, args []string) error {
	fmt.Printf("hcloud %s\n", Version)
	return nil
}

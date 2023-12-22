package version

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/version"
)

func NewCommand(cli *state.State) *cobra.Command {
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
	cmd.Printf("hcloud %s\n", version.Version)
	return nil
}

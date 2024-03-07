package version

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/version"
)

func NewCommand(_ state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "version",
		Short:                 "Print version information",
		Args:                  util.Validate,
		DisableFlagsInUseLine: true,
		RunE:                  runVersion,
	}
	return cmd
}

func runVersion(cmd *cobra.Command, _ []string) error {
	cmd.Printf("hcloud %s\n", version.Version)
	return nil
}

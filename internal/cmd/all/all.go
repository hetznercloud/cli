package all

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "all",
		Short:                 "Commands that apply to all resources",
		Args:                  util.ValidateExact,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		ListCmd.CobraCommand(s),
	)
	return cmd
}

package storagebox

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "storage-box",
		Aliases:               []string{"storage-box"},
		Short:                 "Manage Storage Boxes",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		ListCmd.CobraCommand(s),
		CreateCmd.CobraCommand(s),
		FoldersCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
		LabelCmds.AddCobraCommand(s),
		LabelCmds.RemoveCobraCommand(s),
	)
	return cmd
}

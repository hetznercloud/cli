package volume

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "volume",
		Short:                 "Manage Volumes",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddGroup(&cobra.Group{ID: "general", Title: "General"})
	cmd.AddGroup(&cobra.Group{ID: "protection", Title: "Protection"})
	cmd.AddGroup(&cobra.Group{ID: "attach", Title: "Attach"})

	cmd.AddCommand(
		util.WithGroup("general", ListCmd.CobraCommand(s)),
		util.WithGroup("general", DescribeCmd.CobraCommand(s)),
		util.WithGroup("general", DeleteCmd.CobraCommand(s)),
		util.WithGroup("general", CreateCmd.CobraCommand(s)),
		util.WithGroup("general", UpdateCmd.CobraCommand(s)),
		util.WithGroup("general", LabelCmds.AddCobraCommand(s)),
		util.WithGroup("general", LabelCmds.RemoveCobraCommand(s)),
		util.WithGroup("general", ResizeCmd.CobraCommand(s)),

		util.WithGroup("protection", EnableProtectionCmd.CobraCommand(s)),
		util.WithGroup("protection", DisableProtectionCmd.CobraCommand(s)),

		util.WithGroup("attach", AttachCmd.CobraCommand(s)),
		util.WithGroup("attach", DetachCmd.CobraCommand(s)),
	)
	return cmd
}

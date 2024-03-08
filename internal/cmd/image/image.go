package image

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "image",
		Short:                 "Manage images",
		Args:                  util.ValidateExact,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}

	util.AddGroup(cmd, "general", "General",
		ListCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		UpdateCmd.CobraCommand(s),
		LabelCmds.AddCobraCommand(s),
		LabelCmds.RemoveCobraCommand(s),
	)

	util.AddGroup(cmd, "protection", "Protection",
		EnableProtectionCmd.CobraCommand(s),
		DisableProtectionCmd.CobraCommand(s),
	)
	return cmd
}

package subaccount

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "subaccount",
		Aliases:               []string{"subaccounts"},
		Short:                 "Manage Storage Box Subaccounts",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		ListCmd.CobraCommand(s),
		CreateCmd.CobraCommand(s),
		DeleteCmd.CobraCommand(s),
		UpdateCmd.CobraCommand(s),
		DescribeCmd.CobraCommand(s),
		ResetPasswordCmd.CobraCommand(s),
		UpdateAccessSettingsCmd.CobraCommand(s),
		ChangeHomeDirectoryCmd.CobraCommand(s),
	)
	return experimental.StorageBoxes(s, cmd)
}

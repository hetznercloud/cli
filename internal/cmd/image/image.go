package image

import (
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "image",
		Short:                 "Manage images",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newListCommand(cli),
		newDeleteCommand(cli),
		newDescribeCommand(cli),
		newUpdateCommand(cli),
		newEnableProtectionCommand(cli),
		newDisableProtectionCommand(cli),
		addLabelCmd.CobraCommand(cli.Context, hcapi2.NewClient(cli.Client()), cli),
		removeLabelCmd.CobraCommand(cli.Context, hcapi2.NewClient(cli.Client()), cli),
	)
	return cmd
}

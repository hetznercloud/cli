package iso

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(cli *state.State, client hcapi2.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "iso",
		Short:                 "Manage ISOs",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		ListCmd.CobraCommand(cli.Context, client, cli),
		DescribeCmd.CobraCommand(cli.Context, client, cli),
	)
	return cmd
}

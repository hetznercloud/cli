package certificate

import (
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "certificate",
		Short:                 "Manage certificates",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		listCmd.CobraCommand(cli.Context, hcapi2.NewClient(cli.Client()), cli),
		newCreateCommand(cli),
		newUpdateCommand(cli),
		labelCmds.AddCobraCommand(cli.Context, hcapi2.NewClient(cli.Client()), cli),
		labelCmds.RemoveCobraCommand(cli.Context, hcapi2.NewClient(cli.Client()), cli),
		deleteCmd.CobraCommand(cli.Context, hcapi2.NewClient(cli.Client()), cli),
		describeCmd.CobraCommand(cli.Context, hcapi2.NewClient(cli.Client()), cli),
	)

	return cmd
}

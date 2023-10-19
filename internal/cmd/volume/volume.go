package volume

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

func NewCommand(cli *state.State, client hcapi2.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "volume",
		Short:                 "Manage Volumes",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		ListCmd.CobraCommand(cli.Context, client, cli),
		CreateCommand.CobraCommand(cli.Context, client, cli, cli),
		updateCmd.CobraCommand(cli.Context, client, cli),
		deleteCmd.CobraCommand(cli.Context, client, cli, cli),
		describeCmd.CobraCommand(cli.Context, client, cli),
		AttachCommand.CobraCommand(cli.Context, client, cli, cli),
		DetachCommand.CobraCommand(cli.Context, client, cli, cli),
		ResizeCommand.CobraCommand(cli.Context, client, cli, cli),
		EnableProtectionCommand.CobraCommand(cli.Context, client, cli, cli),
		DisableProtectionCommand.CobraCommand(cli.Context, client, cli, cli),
		labelCmds.AddCobraCommand(cli.Context, client, cli),
		labelCmds.RemoveCobraCommand(cli.Context, client, cli),
	)
	return cmd
}

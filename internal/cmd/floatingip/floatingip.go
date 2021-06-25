package floatingip

import (
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCommand(cli *state.State, client hcapi2.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "floating-ip",
		Short:                 "Manage Floating IPs",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newUpdateCommand(cli),

		listCmd.CobraCommand(cli.Context, client, cli),
		newCreateCommand(cli),
		describeCmd.CobraCommand(cli.Context, client, cli),
		newAssignCommand(cli),
		newUnassignCommand(cli),
		deleteCmd.CobraCommand(cli.Context, client, cli),
		newEnableProtectionCommand(cli),
		newDisableProtectionCommand(cli),
		labelCmds.AddCobraCommand(cli.Context, client, cli),
		labelCmds.RemoveCobraCommand(cli.Context, client, cli),
		newSetRDNSCommand(cli),
	)
	return cmd
}

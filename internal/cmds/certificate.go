package cmds

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func NewCertificatesCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "certificate",
		Short:                 "Manage certificates",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}
	cmd.AddCommand(
		newCertificatesListCommand(cli),
		newCertificateCreateCommand(cli),
		newCertificateUpdateCommand(cli),
		newCertificateAddLabelCommand(cli),
		newCertificateRemoveLabelCommand(cli),
		newCertificateDeleteCommand(cli),
		newCertificateDescribeCommand(cli),
	)

	return cmd
}

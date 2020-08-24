package cli

import "github.com/spf13/cobra"

func newCertificatesCommand(cli *CLI) *cobra.Command {
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

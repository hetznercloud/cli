package cli

import "github.com/spf13/cobra"

func newCertificatesCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "certificate",
		Short:                 "Manage certificates",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  cli.wrap(runCertificates),
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

func runCertificates(cli *CLI, cmd *cobra.Command, args []string) error {
	return cmd.Usage()
}

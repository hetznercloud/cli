package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCertificateDeleteCommand(cli *CLI) *cobra.Command {
	return &cobra.Command{
		Use:                   "delete CERTIFICATE",
		Short:                 "Delete a certificate",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runCertificateDelete),
	}
}

func runCertificateDelete(cli *CLI, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	cert, _, err := cli.Client().Certificate.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if cert == nil {
		return fmt.Errorf("Certificate %s not found", idOrName)
	}
	_, err = cli.Client().Certificate.Delete(cli.Context, cert)
	if err != nil {
		return err
	}
	fmt.Printf("Certificate %d deleted\n", cert.ID)
	return nil
}

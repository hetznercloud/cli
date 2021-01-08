package cmds

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"
)

func newCertificateDeleteCommand(cli *state.State) *cobra.Command {
	return &cobra.Command{
		Use:                   "delete CERTIFICATE",
		Short:                 "Delete a certificate",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.CertificateNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runCertificateDelete),
	}
}

func runCertificateDelete(cli *state.State, cmd *cobra.Command, args []string) error {
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

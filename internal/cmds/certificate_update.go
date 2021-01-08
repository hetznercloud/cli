package cmds

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newCertificateUpdateCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "update [FLAGS] CERTIFICATE",
		Short:                 "Update an existing Certificate",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.CertificateNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runCertificateUpdate),
	}

	cmd.Flags().String("name", "", "Certificate name")
	return cmd
}

func runCertificateUpdate(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	cert, _, err := cli.Client().Certificate.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if cert == nil {
		return fmt.Errorf("Certificate %s not found", idOrName)
	}
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	updOpts := hcloud.CertificateUpdateOpts{
		Name: name,
	}
	_, _, err = cli.Client().Certificate.Update(cli.Context, cert, updOpts)
	if err != nil {
		return err
	}
	fmt.Printf("Certificate %d updated\n", cert.ID)
	return nil
}

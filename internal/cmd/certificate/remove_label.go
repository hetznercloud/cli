package certificate

import (
	"errors"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newRemoveLabelCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-label [FLAGS] CERTIFICATE LABELKEY",
		Short: "Remove a label from a certificate",
		Args:  cobra.RangeArgs(1, 2),
		ValidArgsFunction: cmpl.SuggestArgs(
			cmpl.SuggestCandidatesF(cli.CertificateNames),
			cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
				if len(args) != 1 {
					return nil
				}
				return cli.CertificateLabelKeys(args[0])
			}),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateRemoveLabel, cli.EnsureToken),
		RunE:                  cli.Wrap(runRemoveLabel),
	}
	cmd.Flags().BoolP("all", "a", false, "Remove all labels")
	return cmd
}

func validateRemoveLabel(cmd *cobra.Command, args []string) error {
	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		return err
	}
	if all && len(args) != 1 {
		return errors.New("must not specify a label key when using --all/-a")
	}
	if !all && len(args) != 2 {
		return errors.New("must specify a label key when not using --all/-a")
	}
	return nil
}

func runRemoveLabel(cli *state.State, cmd *cobra.Command, args []string) error {
	// We ensured the all flag is a valid boolean in
	// validateCertificateRemoveLabel. No need to handle the error again here.
	all, _ := cmd.Flags().GetBool("all")
	idOrName := args[0]
	cert, _, err := cli.Client().Certificate.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if cert == nil {
		return fmt.Errorf("Certificate not found: %s", idOrName)
	}
	if all {
		cert.Labels = make(map[string]string)
	} else {
		label := args[1]
		if _, ok := cert.Labels[label]; !ok {
			return fmt.Errorf("Label %s on certificate %d does not exist", label, cert.ID)
		}
		delete(cert.Labels, label)
	}
	opts := hcloud.CertificateUpdateOpts{
		Labels: cert.Labels,
	}
	_, _, err = cli.Client().Certificate.Update(cli.Context, cert, opts)
	if err != nil {
		return err
	}

	if all {
		fmt.Printf("All labels removed from certificate %d\n", cert.ID)
	} else {
		fmt.Printf("Label %s removed from certificate %d\n", args[1], cert.ID)
	}
	return nil
}

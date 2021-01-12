package certificate

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newAddLabelCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add-label [FLAGS] CERTIFICATE LABEL",
		Short:                 "Add a label to a certificate",
		Args:                  cobra.ExactArgs(2),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.CertificateNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateAddLabel, cli.EnsureToken),
		RunE:                  cli.Wrap(runAddLabel),
	}

	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite label if it exists already")
	return cmd
}

func validateAddLabel(cmd *cobra.Command, args []string) error {
	label := util.SplitLabel(args[1])
	if len(label) != 2 {
		return fmt.Errorf("invalid label: %s", args[1])
	}
	return nil
}

func runAddLabel(cli *state.State, cmd *cobra.Command, args []string) error {
	overwrite, err := cmd.Flags().GetBool("overwrite")
	if err != nil {
		return err
	}
	idOrName := args[0]
	cert, _, err := cli.Client().Certificate.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if cert == nil {
		return fmt.Errorf("Certificate not found: %s", idOrName)
	}
	label := util.SplitLabel(args[1])
	if _, ok := cert.Labels[label[0]]; ok && !overwrite {
		return fmt.Errorf("Label %s on certificate %d already exists", label[0], cert.ID)
	}
	if cert.Labels == nil {
		cert.Labels = make(map[string]string)
	}
	labels := cert.Labels
	labels[label[0]] = label[1]
	opts := hcloud.CertificateUpdateOpts{
		Labels: labels,
	}
	_, _, err = cli.Client().Certificate.Update(cli.Context, cert, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Label %s added to certificate %d\n", label[0], cert.ID)
	return nil
}

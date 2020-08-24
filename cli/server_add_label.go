package cli

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newServerAddLabelCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add-label [FLAGS] SERVER LABEL",
		Short:                 "Add a label to a server",
		Args:                  cobra.ExactArgs(2),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.ServerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               chainRunE(validateServerAddLabel, cli.ensureToken),
		RunE:                  cli.wrap(runServerAddLabel),
	}

	cmd.Flags().BoolP("overwrite", "o", false, "Overwrite label if it exists already")
	return cmd
}

func validateServerAddLabel(cmd *cobra.Command, args []string) error {
	label := splitLabel(args[1])
	if len(label) != 2 {
		return fmt.Errorf("invalid label: %s", args[1])
	}

	return nil
}

func runServerAddLabel(cli *CLI, cmd *cobra.Command, args []string) error {
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}
	label := splitLabel(args[1])

	if _, ok := server.Labels[label[0]]; ok && !overwrite {
		return fmt.Errorf("label %s on server %d already exists", label[0], server.ID)
	}
	labels := server.Labels
	labels[label[0]] = label[1]
	opts := hcloud.ServerUpdateOpts{
		Labels: labels,
	}
	_, _, err = cli.Client().Server.Update(cli.Context, server, opts)
	if err != nil {
		return err
	}
	fmt.Printf("Label %s added to server %d\n", label[0], server.ID)

	return nil
}

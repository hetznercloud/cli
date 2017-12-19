package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newSSHKeyListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "list",
		Short:            "List SSH keys",
		TraverseChildren: true,
		RunE:             cli.wrap(runSSHKeyList),
	}
	return cmd
}

func runSSHKeyList(cli *CLI, cmd *cobra.Command, args []string) error {
	sshKeys, err := cli.Client().SSHKey.All(cli.Context)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tFINGERPRINT")
	for _, sshKey := range sshKeys {
		fmt.Fprintf(w, "%d\t%.50s\t%s\n", sshKey.ID, sshKey.Name, sshKey.Fingerprint)
	}
	w.Flush()

	return nil
}

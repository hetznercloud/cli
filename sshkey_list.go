package cli

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var sshKeyListTableOutput *tableOutput

func init() {
	sshKeyListTableOutput = newTableOutput().
		AddAllowedFields(hcloud.SSHKey{})
}

func newSSHKeyListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List SSH keys",
		Long: listLongDescription(
			"Displays a list of SSH keys.",
			sshKeyListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runSSHKeyList),
	}
	addListOutputFlag(cmd, sshKeyListTableOutput.Columns())
	return cmd
}

func runSSHKeyList(cli *CLI, cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetStringArray("output")
	outOpts, err := parseOutputOpts(out)
	if err != nil {
		return err
	}

	sshKeys, err := cli.Client().SSHKey.All(cli.Context)
	if err != nil {
		return err
	}

	cols := []string{"id", "name", "fingerprint"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := sshKeyListTableOutput
	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, sshKey := range sshKeys {
		tw.Write(cols, sshKey)
	}
	tw.Flush()
	return nil
}

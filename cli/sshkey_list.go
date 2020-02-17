package cli

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var sshKeyListTableOutput *tableOutput

func init() {
	sshKeyListTableOutput = newTableOutput().
		AddAllowedFields(hcloud.SSHKey{}).
		AddFieldOutputFn("labels", fieldOutputFn(func(obj interface{}) string {
			sshKey := obj.(*hcloud.SSHKey)
			return labelsToString(sshKey.Labels)
		}))
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
	addOutputFlag(cmd, outputOptionNoHeader(), outputOptionColumns(sshKeyListTableOutput.Columns()), outputOptionJSON())
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	return cmd
}

func runSSHKeyList(cli *CLI, cmd *cobra.Command, args []string) error {
	outOpts := outputFlagsForCommand(cmd)

	labelSelector, _ := cmd.Flags().GetString("selector")
	opts := hcloud.SSHKeyListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: labelSelector,
			PerPage:       50,
		},
	}
	sshKeys, err := cli.Client().SSHKey.AllWithOpts(cli.Context, opts)
	if err != nil {
		return err
	}

	if outOpts.IsSet("json") {
		describeJSON(sshKeys, false)
		return nil
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

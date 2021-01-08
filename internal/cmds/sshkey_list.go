package cmds

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"
)

var sshKeyListTableOutput *tableOutput

func init() {
	sshKeyListTableOutput = newTableOutput().
		AddAllowedFields(hcloud.SSHKey{}).
		AddFieldOutputFn("labels", fieldOutputFn(func(obj interface{}) string {
			sshKey := obj.(*hcloud.SSHKey)
			return labelsToString(sshKey.Labels)
		})).
		AddFieldOutputFn("created", fieldOutputFn(func(obj interface{}) string {
			sshKey := obj.(*hcloud.SSHKey)
			return datetime(sshKey.Created)
		}))
}

func newSSHKeyListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List SSH keys",
		Long: listLongDescription(
			"Displays a list of SSH keys.",
			sshKeyListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runSSHKeyList),
	}
	addOutputFlag(cmd, outputOptionNoHeader(), outputOptionColumns(sshKeyListTableOutput.Columns()), outputOptionJSON())
	cmd.Flags().StringP("selector", "l", "", "Selector to filter by labels")
	return cmd
}

func runSSHKeyList(cli *state.State, cmd *cobra.Command, args []string) error {
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
		var sshKeySchemas []schema.SSHKey
		for _, sshKey := range sshKeys {
			sshKeySchema := schema.SSHKey{
				ID:          sshKey.ID,
				Name:        sshKey.Name,
				Fingerprint: sshKey.Fingerprint,
				PublicKey:   sshKey.PublicKey,
				Labels:      sshKey.Labels,
				Created:     sshKey.Created,
			}
			sshKeySchemas = append(sshKeySchemas, sshKeySchema)
		}
		return describeJSON(sshKeySchemas)
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

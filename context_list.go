package cli

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newContextListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list [FLAGS]",
		Short:                 "List contexts",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE: cli.wrap(runContextList),
	}
	return cmd
}

func runContextList(cli *CLI, cmd *cobra.Command, args []string) error {
	out, _ := cmd.Flags().GetStringArray("output")
	outOpts, err := parseOutputOpts(out)
	if err != nil {
		return err
	}

	cols := []string{"name"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := newTableOutput().
		AddAllowedFields(hcloud.Datacenter{}).
		RemoveAllowedField("token")

	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, context := range cli.Config.Contexts {
		tw.Write(cols, context)
	}
	tw.Flush()
	return nil
}

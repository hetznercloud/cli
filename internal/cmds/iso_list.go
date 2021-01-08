package cmds

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"
)

var isoListTableOutput *tableOutput

func init() {
	isoListTableOutput = newTableOutput().
		AddAllowedFields(hcloud.ISO{})
}

func newISOListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List ISOs",
		Long: listLongDescription(
			"Displays a list of ISOs.",
			isoListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runISOList),
	}
	addOutputFlag(cmd, outputOptionNoHeader(), outputOptionColumns(isoListTableOutput.Columns()), outputOptionJSON())
	return cmd
}

func runISOList(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := outputFlagsForCommand(cmd)

	isos, err := cli.Client().ISO.All(cli.Context)
	if err != nil {
		return err
	}

	if outOpts.IsSet("json") {
		var isoSchemas []schema.ISO
		for _, iso := range isos {
			isoSchemas = append(isoSchemas, isoToSchema(*iso))
		}
		return describeJSON(isoSchemas)
	}

	cols := []string{"id", "name", "description", "type"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := isoListTableOutput
	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, iso := range isos {
		tw.Write(cols, iso)
	}
	tw.Flush()

	return nil
}

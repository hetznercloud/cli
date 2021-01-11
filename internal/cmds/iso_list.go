package cmds

import (
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"
)

var isoListTableOutput *output.Table

func init() {
	isoListTableOutput = output.NewTable().
		AddAllowedFields(hcloud.ISO{})
}

func newISOListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List ISOs",
		Long: util.ListLongDescription(
			"Displays a list of ISOs.",
			isoListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runISOList),
	}
	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(isoListTableOutput.Columns()), output.OptionJSON())
	return cmd
}

func runISOList(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := output.FlagsForCommand(cmd)

	isos, err := cli.Client().ISO.All(cli.Context)
	if err != nil {
		return err
	}

	if outOpts.IsSet("json") {
		var isoSchemas []schema.ISO
		for _, iso := range isos {
			isoSchemas = append(isoSchemas, util.ISOToSchema(*iso))
		}
		return util.DescribeJSON(isoSchemas)
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

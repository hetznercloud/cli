package cmds

import (
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"
)

var locationListTableOutput *output.Table

func init() {
	locationListTableOutput = output.NewTable().
		AddAllowedFields(hcloud.Location{})
}

func newLocationListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List locations",
		Long: util.ListLongDescription(
			"Displays a list of locations.",
			locationListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runLocationList),
	}
	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(locationListTableOutput.Columns()), output.OptionJSON())
	return cmd
}

func runLocationList(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := output.FlagsForCommand(cmd)

	locations, err := cli.Client().Location.All(cli.Context)
	if err != nil {
		return err
	}

	if outOpts.IsSet("json") {
		var locationSchemas []schema.Location
		for _, location := range locations {
			locationSchemas = append(locationSchemas, util.LocationToSchema(*location))
		}
		return util.DescribeJSON(locationSchemas)
	}

	cols := []string{"id", "name", "description", "network_zone", "country", "city"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := locationListTableOutput
	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, location := range locations {
		tw.Write(cols, location)
	}
	tw.Flush()

	return nil
}

package cli

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var locationListTableOutput *tableOutput

func init() {
	locationListTableOutput = newTableOutput().
		AddAllowedFields(hcloud.Location{})
}

func newLocationListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List locations",
		Long: listLongDescription(
			"Displays a list of locations.",
			locationListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runLocationList),
	}
	addOutputFlag(cmd, outputOptionNoHeader(), outputOptionColumns(locationListTableOutput.Columns()), outputOptionJSON())
	return cmd
}

func runLocationList(cli *CLI, cmd *cobra.Command, args []string) error {
	outOpts := outputFlagsForCommand(cmd)

	locations, err := cli.Client().Location.All(cli.Context)
	if err != nil {
		return err
	}

	if outOpts.IsSet("json") {
		describeJSON(locations, false)
		return nil
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

package cmds

import (
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"
)

var datacenterListTableOutput *tableOutput

func init() {
	datacenterListTableOutput = newTableOutput().
		AddAllowedFields(hcloud.Datacenter{}).
		AddFieldOutputFn("location", fieldOutputFn(func(obj interface{}) string {
			datacenter := obj.(*hcloud.Datacenter)
			return datacenter.Location.Name
		}))
}

func newDatacenterListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List datacenters",
		Long: listLongDescription(
			"Displays a list of datacenters.",
			datacenterListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runDatacenterList),
	}
	addOutputFlag(cmd, outputOptionNoHeader(), outputOptionColumns(datacenterListTableOutput.Columns()), outputOptionJSON())
	return cmd
}

func runDatacenterList(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := outputFlagsForCommand(cmd)

	datacenters, err := cli.Client().Datacenter.All(cli.Context)

	if outOpts.IsSet("json") {
		var datacenterSchemas []schema.Datacenter
		for _, datacenter := range datacenters {
			datacenterSchemas = append(datacenterSchemas, datacenterToSchema(*datacenter))
		}
		return describeJSON(datacenterSchemas)
	}

	if err != nil {
		return err
	}

	cols := []string{"id", "name", "description", "location"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := datacenterListTableOutput
	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, datacenter := range datacenters {
		tw.Write(cols, datacenter)
	}
	tw.Flush()
	return nil
}

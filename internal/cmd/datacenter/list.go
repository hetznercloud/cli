package datacenter

import (
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"
)

var datacenterListTableOutput *output.Table

func init() {
	datacenterListTableOutput = output.NewTable().
		AddAllowedFields(hcloud.Datacenter{}).
		AddFieldFn("location", output.FieldFn(func(obj interface{}) string {
			datacenter := obj.(*hcloud.Datacenter)
			return datacenter.Location.Name
		}))
}

func newListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List datacenters",
		Long: util.ListLongDescription(
			"Displays a list of datacenters.",
			datacenterListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runDatacenterList),
	}
	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(datacenterListTableOutput.Columns()), output.OptionJSON())
	return cmd
}

func runDatacenterList(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := output.FlagsForCommand(cmd)

	datacenters, err := cli.Client().Datacenter.All(cli.Context)

	if outOpts.IsSet("json") {
		var datacenterSchemas []schema.Datacenter
		for _, datacenter := range datacenters {
			datacenterSchemas = append(datacenterSchemas, util.DatacenterToSchema(*datacenter))
		}
		return util.DescribeJSON(datacenterSchemas)
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

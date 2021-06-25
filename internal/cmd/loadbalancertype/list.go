package loadbalancertype

import (
	"context"

	"github.com/hetznercloud/cli/internal/state"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

var ListCmd = base.ListCmd{
	ResourceNamePlural: "Load Balancer Types",

	DefaultColumns: []string{"id", "name", "description", "max_services", "max_connections", "max_targets"},

	Fetch: func(ctx context.Context, client hcapi2.Client, listOpts hcloud.ListOpts) ([]interface{}, error) {
		loadBalancerTypes, _, err := client.LoadBalancerType().List(ctx, hcloud.LoadBalancerTypeListOpts{ListOpts: listOpts})

		var resources []interface{}
		for _, r := range loadBalancerTypes {
			resources = append(resources, r)
		}
		return resources, err
	},

	OutputTable: func(client hcapi2.Client) *output.Table {
		return output.NewTable().
			AddAllowedFields(hcloud.LoadBalancerType{})
	},

	JSONSchema: func(resources []interface{}) interface{} {
		var loadBalancerTypeSchemas []schema.LoadBalancerType
		for _, resource := range resources {
			loadBalancerType := resource.(*hcloud.LoadBalancerType)
			loadBalancerTypeSchemas = append(loadBalancerTypeSchemas, util.LoadBalancerTypeToSchema(*loadBalancerType))
		}
		return loadBalancerTypeSchemas
	},
}

var listTableOutput *output.Table

func init() {
	listTableOutput = output.NewTable().
		AddAllowedFields(hcloud.LoadBalancerType{})
}

func newListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List Load Balancer types",
		Long: util.ListLongDescription(
			"Displays a list of Load Balancer types.",
			listTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runList),
	}
	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(listTableOutput.Columns()), output.OptionJSON())
	return cmd
}

func runList(cli *state.State, cmd *cobra.Command, args []string) error {
	outOpts := output.FlagsForCommand(cmd)

	loadBalancerTypes, err := cli.Client().LoadBalancerType.All(cli.Context)
	if err != nil {
		return err
	}

	if outOpts.IsSet("json") {
		var loadBalancerTypeSchemas []schema.LoadBalancerType
		for _, loadBalancerType := range loadBalancerTypes {
			loadBalancerTypeSchemas = append(loadBalancerTypeSchemas, util.LoadBalancerTypeToSchema(*loadBalancerType))
		}
		return util.DescribeJSON(loadBalancerTypeSchemas)
	}

	cols := []string{"id", "name", "description", "max_services", "max_connections", "max_targets"}
	if outOpts.IsSet("columns") {
		cols = outOpts["columns"]
	}

	tw := listTableOutput
	if err = tw.ValidateColumns(cols); err != nil {
		return err
	}

	if !outOpts.IsSet("noheader") {
		tw.WriteHeader(cols)
	}
	for _, loadBalancerType := range loadBalancerTypes {
		tw.Write(cols, loadBalancerType)
	}
	tw.Flush()

	return nil
}

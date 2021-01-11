package cmds

import (
	"github.com/hetznercloud/cli/internal/cmd/output"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"
)

var loadBalancerTypeListTableOutput *output.Table

func init() {
	loadBalancerTypeListTableOutput = output.NewTable().
		AddAllowedFields(hcloud.LoadBalancerType{})
}

func newLoadBalancerTypeListCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List Load Balancer types",
		Long: util.ListLongDescription(
			"Displays a list of Load Balancer types.",
			loadBalancerTypeListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runLoadBalancerTypeList),
	}
	output.AddFlag(cmd, output.OptionNoHeader(), output.OptionColumns(loadBalancerTypeListTableOutput.Columns()), output.OptionJSON())
	return cmd
}

func runLoadBalancerTypeList(cli *state.State, cmd *cobra.Command, args []string) error {
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

	tw := loadBalancerTypeListTableOutput
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

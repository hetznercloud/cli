package cli

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/hetznercloud/hcloud-go/hcloud/schema"
	"github.com/spf13/cobra"
)

var loadBalancerTypeListTableOutput *tableOutput

func init() {
	loadBalancerTypeListTableOutput = newTableOutput().
		AddAllowedFields(hcloud.LoadBalancerType{})
}

func newLoadBalancerTypeListCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [FLAGS]",
		Short: "List Load Balancer types",
		Long: listLongDescription(
			"Displays a list of Load Balancer types.",
			loadBalancerTypeListTableOutput.Columns(),
		),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runLoadBalancerTypeList),
	}
	addOutputFlag(cmd, outputOptionNoHeader(), outputOptionColumns(loadBalancerTypeListTableOutput.Columns()), outputOptionJSON())
	return cmd
}

func runLoadBalancerTypeList(cli *CLI, cmd *cobra.Command, args []string) error {
	outOpts := outputFlagsForCommand(cmd)

	loadBalancerTypes, err := cli.Client().LoadBalancerType.All(cli.Context)
	if err != nil {
		return err
	}

	if outOpts.IsSet("json") {
		var loadBalancerTypeSchemas []schema.LoadBalancerType
		for _, loadBalancerType := range loadBalancerTypes {
			loadBalancerTypeSchemas = append(loadBalancerTypeSchemas, loadBalancerTypeToSchema(*loadBalancerType))
		}
		return describeJSON(loadBalancerTypeSchemas)
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

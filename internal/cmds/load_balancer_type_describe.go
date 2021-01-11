package cmds

import (
	"encoding/json"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newLoadBalancerTypenDescribeCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "describe [FLAGS] LOADBALANCERTYPE",
		Short:                 "Describe a Load Balancer type",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.LoadBalancerTypeNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.EnsureToken,
		RunE:                  cli.Wrap(runLoadBalancerTypeDescribe),
	}
	addOutputFlag(cmd, outputOptionJSON(), outputOptionFormat())
	return cmd
}

func runLoadBalancerTypeDescribe(cli *state.State, cmd *cobra.Command, args []string) error {
	outputFlags := outputFlagsForCommand(cmd)

	idOrName := args[0]
	loadBalancerType, resp, err := cli.Client().LoadBalancerType.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancerType == nil {
		return fmt.Errorf("loadBalancerType not found: %s", idOrName)
	}

	switch {
	case outputFlags.IsSet("json"):
		return loadBalancerTypeDescribeJSON(resp)
	case outputFlags.IsSet("format"):
		return util.DescribeFormat(loadBalancerType, outputFlags["format"][0])
	default:
		return loadBalancerTypeDescribeText(cli, loadBalancerType)
	}
}

func loadBalancerTypeDescribeText(cli *state.State, loadBalancerType *hcloud.LoadBalancerType) error {
	fmt.Printf("ID:\t\t\t\t%d\n", loadBalancerType.ID)
	fmt.Printf("Name:\t\t\t\t%s\n", loadBalancerType.Name)
	fmt.Printf("Description:\t\t\t%s\n", loadBalancerType.Description)
	fmt.Printf("Max Services:\t\t\t%d\n", loadBalancerType.MaxServices)
	fmt.Printf("Max Connections:\t\t%d\n", loadBalancerType.MaxConnections)
	fmt.Printf("Max Targets:\t\t\t%d\n", loadBalancerType.MaxTargets)
	fmt.Printf("Max assigned Certificates:\t%d\n", loadBalancerType.MaxAssignedCertificates)

	fmt.Printf("Pricings per Location:\n")
	for _, price := range loadBalancerType.Pricings {
		fmt.Printf("  - Location:\t%s:\n", price.Location.Name)
		fmt.Printf("    Hourly:\t€ %s\n", price.Hourly.Gross)
		fmt.Printf("    Monthly:\t€ %s\n", price.Monthly.Gross)
	}
	return nil
}

func loadBalancerTypeDescribeJSON(resp *hcloud.Response) error {
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	if loadBalancerType, ok := data["loadBalancerType"]; ok {
		return util.DescribeJSON(loadBalancerType)
	}
	if loadBalancerTypes, ok := data["loadBalancerTypes"].([]interface{}); ok {
		return util.DescribeJSON(loadBalancerTypes[0])
	}
	return util.DescribeJSON(data)
}

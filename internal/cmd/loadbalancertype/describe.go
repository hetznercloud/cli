package loadbalancertype

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd{
	ResourceNameSingular: "Load Balancer Type",
	ShortDescription:     "Describe a Load Balancer type",
	JSONKeyGetByID:       "load_balancer_type",
	JSONKeyGetByName:     "load_balancer_types",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.LoadBalancerType().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, interface{}, error) {
		lbt, _, err := s.LoadBalancerType().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return lbt, hcloud.SchemaFromLoadBalancerType(lbt), nil
	},
	PrintText: func(_ state.State, cmd *cobra.Command, resource interface{}) error {
		loadBalancerType := resource.(*hcloud.LoadBalancerType)

		cmd.Printf("ID:\t\t\t\t%d\n", loadBalancerType.ID)
		cmd.Printf("Name:\t\t\t\t%s\n", loadBalancerType.Name)
		cmd.Printf("Description:\t\t\t%s\n", loadBalancerType.Description)
		cmd.Printf("Max Services:\t\t\t%d\n", loadBalancerType.MaxServices)
		cmd.Printf("Max Connections:\t\t%d\n", loadBalancerType.MaxConnections)
		cmd.Printf("Max Targets:\t\t\t%d\n", loadBalancerType.MaxTargets)
		cmd.Printf("Max assigned Certificates:\t%d\n", loadBalancerType.MaxAssignedCertificates)

		cmd.Printf("Pricings per Location:\n")
		for _, price := range loadBalancerType.Pricings {
			cmd.Printf("  - Location:\t%s:\n", price.Location.Name)
			cmd.Printf("    Hourly:\t€ %s\n", price.Hourly.Gross)
			cmd.Printf("    Monthly:\t€ %s\n", price.Monthly.Gross)
		}
		return nil
	},
}

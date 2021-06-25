package loadbalancertype

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

var describeCmd = base.DescribeCmd{
	ResourceNameSingular: "Load Balancer Type",
	ShortDescription:     "Describe a Load Balancer type",
	JSONKeyGetByID:       "load_balancer_type",
	JSONKeyGetByName:     "load_balancer_types",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.LoadBalancerType().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.LoadBalancerType().Get(ctx, idOrName)
	},
	PrintText: func(_ context.Context, _ hcapi2.Client, _ *cobra.Command, resource interface{}) error {
		loadBalancerType := resource.(*hcloud.LoadBalancerType)

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
	},
}

package loadbalancer

import (
	"context"
	"net"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var setRDNSCmd = base.SetRdnsCmd{
	ResourceNameSingular: "Load Balancer",
	ShortDescription:     "Change reverse DNS of a Load Balancer",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.LoadBalancer().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.LoadBalancer().Get(ctx, idOrName)
	},
	GetDefaultIP: func(resource interface{}) net.IP {
		loadBalancer := resource.(*hcloud.LoadBalancer)
		return loadBalancer.PublicNet.IPv4.IP
	},
}

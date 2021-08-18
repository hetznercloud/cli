package loadbalancer

import (
	"context"
	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"net"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var setRDNSCmd = base.SetRdnsCmd{
	ResourceNameSingular: "Load Balancer",
	ShortDescription:     "Change reverse DNS of a Load Balancer",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.FloatingIP().Names },
	Fetch: func(ctx context.Context, client hcapi2.Client, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return client.LoadBalancer().Get(ctx, idOrName)
	},
	GetDefaultIP: func(resource interface{}) net.IP {
		loadBalancer := resource.(*hcloud.LoadBalancer)
		return loadBalancer.PublicNet.IPv4.IP
	},
}

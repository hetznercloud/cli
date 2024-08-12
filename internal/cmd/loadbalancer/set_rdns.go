package loadbalancer

import (
	"net"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var SetRDNSCmd = base.SetRdnsCmd{
	ResourceNameSingular: "Load Balancer",
	ShortDescription:     "Change reverse DNS of a Load Balancer",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.LoadBalancer().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.Client().LoadBalancer().Get(s, idOrName)
	},
	GetDefaultIP: func(resource interface{}) net.IP {
		loadBalancer := resource.(*hcloud.LoadBalancer)
		return loadBalancer.PublicNet.IPv4.IP
	},
}

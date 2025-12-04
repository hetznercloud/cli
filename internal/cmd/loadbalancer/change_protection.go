package loadbalancer

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ChangeProtectionCmds = base.ChangeProtectionCmds[*hcloud.LoadBalancer, hcloud.LoadBalancerChangeProtectionOpts]{
	ResourceNameSingular: "Load Balancer",

	NameSuggestions: func(client hcapi2.Client) func() []string {
		return client.LoadBalancer().Names
	},

	ProtectionLevels: map[string]func(opts *hcloud.LoadBalancerChangeProtectionOpts, value bool){
		"delete": func(opts *hcloud.LoadBalancerChangeProtectionOpts, value bool) {
			opts.Delete = &value
		},
	},

	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.LoadBalancer, *hcloud.Response, error) {
		return s.Client().LoadBalancer().Get(s, idOrName)
	},

	ChangeProtectionFunction: func(s state.State, loadBalancer *hcloud.LoadBalancer, opts hcloud.LoadBalancerChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error) {
		return s.Client().LoadBalancer().ChangeProtection(s, loadBalancer, opts)
	},

	IDOrName: func(loadBalancer *hcloud.LoadBalancer) string {
		return fmt.Sprint(loadBalancer.ID)
	},
}

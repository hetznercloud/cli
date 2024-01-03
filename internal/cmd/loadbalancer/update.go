package loadbalancer

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var UpdateCmd = base.UpdateCmd{
	ResourceNameSingular: "Load Balancer",
	ShortDescription:     "Update a Load Balancer",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.LoadBalancer().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.LoadBalancer().Get(s, idOrName)
	},
	DefineFlags: func(cmd *cobra.Command) {
		cmd.Flags().String("name", "", "LoadBalancer name")
	},
	Update: func(s state.State, cmd *cobra.Command, resource interface{}, flags map[string]pflag.Value) error {
		floatingIP := resource.(*hcloud.LoadBalancer)
		updOpts := hcloud.LoadBalancerUpdateOpts{
			Name: flags["name"].String(),
		}
		_, _, err := s.LoadBalancer().Update(s, floatingIP, updOpts)
		if err != nil {
			return err
		}
		return nil
	},
}

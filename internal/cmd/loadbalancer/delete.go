package loadbalancer

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd{
	ResourceNameSingular: "Load Balancer",
	ShortDescription:     "Delete a Load Balancer",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.LoadBalancer().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (interface{}, *hcloud.Response, error) {
		return s.Client().LoadBalancer().Get(s, idOrName)
	},
	Delete: func(s state.State, cmd *cobra.Command, resource interface{}) error {
		loadBalancer := resource.(*hcloud.LoadBalancer)
		if _, err := s.Client().LoadBalancer().Delete(s, loadBalancer); err != nil {
			return err
		}
		return nil
	},
}

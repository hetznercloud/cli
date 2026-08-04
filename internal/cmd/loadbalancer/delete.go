package loadbalancer

import (
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DeleteCmd = base.DeleteCmd[*hcloud.LoadBalancer]{
	ResourceNameSingular: "Load Balancer",
	ResourceNamePlural:   "Load Balancers",
	ShortDescription:     "Delete a Load Balancer",
	NameSuggestions:      func(c hcapi2.Client) hcapi2.CompletionFunc { return c.LoadBalancer().Names },
	Fetch: func(s state.State, cmd *cobra.Command, idOrName string) (*hcloud.LoadBalancer, *hcloud.Response, error) {
		return s.Client().LoadBalancer().Get(cmd.Context(), idOrName)
	},
	Delete: func(s state.State, cmd *cobra.Command, loadBalancer *hcloud.LoadBalancer) ([]*hcloud.Action, error) {
		_, err := s.Client().LoadBalancer().Delete(cmd.Context(), loadBalancer)
		return nil, err
	},
}

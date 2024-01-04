package loadbalancer

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds{
	ResourceNameSingular:   "Load Balancer",
	ShortDescriptionAdd:    "Add a label to a Load Balancer",
	ShortDescriptionRemove: "Remove a label from a Load Balancer",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.LoadBalancer().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.LoadBalancer().LabelKeys },
	FetchLabels: func(s state.State, idOrName string) (map[string]string, int64, error) {
		loadBalancer, _, err := s.Client().LoadBalancer().Get(s, idOrName)
		if err != nil {
			return nil, 0, err
		}
		if loadBalancer == nil {
			return nil, 0, fmt.Errorf("load balancer not found: %s", idOrName)
		}
		return loadBalancer.Labels, loadBalancer.ID, nil
	},
	SetLabels: func(s state.State, id int64, labels map[string]string) error {
		opts := hcloud.LoadBalancerUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().LoadBalancer().Update(s, &hcloud.LoadBalancer{ID: id}, opts)
		return err
	},
}

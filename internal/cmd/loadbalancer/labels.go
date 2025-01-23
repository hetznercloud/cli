package loadbalancer

import (
	"fmt"
	"strconv"

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
	Fetch: func(s state.State, idOrName string) (any, error) {
		loadBalancer, _, err := s.Client().LoadBalancer().Get(s, idOrName)
		if err != nil {
			return nil, err
		}
		if loadBalancer == nil {
			return nil, fmt.Errorf("load balancer not found: %s", idOrName)
		}
		return loadBalancer, nil
	},
	SetLabels: func(s state.State, resource any, labels map[string]string) error {
		loadBalancer := resource.(*hcloud.LoadBalancer)
		opts := hcloud.LoadBalancerUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().LoadBalancer().Update(s, loadBalancer, opts)
		return err
	},
	GetLabels: func(resource any) map[string]string {
		loadBalancer := resource.(*hcloud.LoadBalancer)
		return loadBalancer.Labels
	},
	GetIDOrName: func(resource any) string {
		loadBalancer := resource.(*hcloud.LoadBalancer)
		return strconv.FormatInt(loadBalancer.ID, 10)
	},
}

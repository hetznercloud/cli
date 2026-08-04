package loadbalancer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds[*hcloud.LoadBalancer]{
	ResourceNameSingular:   "Load Balancer",
	ShortDescriptionAdd:    "Add a label to a Load Balancer",
	ShortDescriptionRemove: "Remove a label from a Load Balancer",
	NameSuggestions:        func(c hcapi2.Client) hcapi2.CompletionFunc { return c.LoadBalancer().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) hcapi2.LabelCompletionFunc { return c.LoadBalancer().LabelKeys },
	Fetch: func(ctx context.Context, s state.State, idOrName string) (*hcloud.LoadBalancer, error) {
		loadBalancer, _, err := s.Client().LoadBalancer().Get(ctx, idOrName)
		if err != nil {
			return nil, err
		}
		if loadBalancer == nil {
			return nil, fmt.Errorf("Load Balancer not found: %s", idOrName)
		}
		return loadBalancer, nil
	},
	SetLabels: func(ctx context.Context, s state.State, loadBalancer *hcloud.LoadBalancer, labels map[string]string) error {
		opts := hcloud.LoadBalancerUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().LoadBalancer().Update(ctx, loadBalancer, opts)
		return err
	},
	GetLabels: func(loadBalancer *hcloud.LoadBalancer) map[string]string {
		return loadBalancer.Labels
	},
	GetIDOrName: func(loadBalancer *hcloud.LoadBalancer) string {
		return strconv.FormatInt(loadBalancer.ID, 10)
	},
}

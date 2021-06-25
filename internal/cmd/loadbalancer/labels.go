package loadbalancer

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

var labelCmds = base.LabelCmds{
	ResourceNameSingular:   "Load Balancer",
	ShortDescriptionAdd:    "Add a label to a Load Balancer",
	ShortDescriptionRemove: "Remove a label from a Load Balancer",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.LoadBalancer().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.LoadBalancer().LabelKeys },
	FetchLabels: func(ctx context.Context, client hcapi2.Client, idOrName string) (map[string]string, int, error) {
		loadBalancer, _, err := client.LoadBalancer().Get(ctx, idOrName)
		if err != nil {
			return nil, 0, err
		}
		if loadBalancer == nil {
			return nil, 0, fmt.Errorf("load balancer not found: %s", idOrName)
		}
		return loadBalancer.Labels, loadBalancer.ID, nil
	},
	SetLabels: func(ctx context.Context, client hcapi2.Client, id int, labels map[string]string) error {
		opts := hcloud.LoadBalancerUpdateOpts{
			Labels: labels,
		}
		_, _, err := client.LoadBalancer().Update(ctx, &hcloud.LoadBalancer{ID: id}, opts)
		return err
	},
}

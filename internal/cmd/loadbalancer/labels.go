package loadbalancer

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var LabelCmds = base.LabelCmds[*hcloud.LoadBalancer]{
	ResourceNameSingular:   "Load Balancer",
	ShortDescriptionAdd:    "Add a label to a Load Balancer",
	ShortDescriptionRemove: "Remove a label from a Load Balancer",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.LoadBalancer().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.LoadBalancer().LabelKeys },
	Fetch: func(s state.State, idOrName string) (*hcloud.LoadBalancer, error) {
		loadBalancer, _, err := s.Client().LoadBalancer().Get(s, idOrName)
		if err != nil {
			return nil, err
		}
		if loadBalancer == nil {
			return nil, fmt.Errorf("Load Balancer not found: %s", idOrName)
		}
		return loadBalancer, nil
	},
	SetLabels: func(s state.State, loadBalancer *hcloud.LoadBalancer, labels map[string]string) error {
		opts := hcloud.LoadBalancerUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().LoadBalancer().Update(s, loadBalancer, opts)
		return err
	},
	GetLabels: func(loadBalancer *hcloud.LoadBalancer) map[string]string {
		return loadBalancer.Labels
	},
	GetIDOrName: func(loadBalancer *hcloud.LoadBalancer) string {
		return strconv.FormatInt(loadBalancer.ID, 10)
	},
	FetchBatch: func(s state.State, idOrNames []string) ([]*hcloud.LoadBalancer, []error) {
		loadBalancers := make([]*hcloud.LoadBalancer, len(idOrNames))
		errors := make([]error, len(idOrNames))

		var wg sync.WaitGroup
		for i, idOrName := range idOrNames {
			wg.Add(1)
			go func(idx int, id string) {
				defer wg.Done()
				loadBalancer, _, err := s.Client().LoadBalancer().Get(s, id)
				if err != nil {
					errors[idx] = err
					return
				}
				if loadBalancer == nil {
					errors[idx] = fmt.Errorf("Load Balancer not found: %s", id)
					return
				}
				loadBalancers[idx] = loadBalancer
			}(i, idOrName)
		}
		wg.Wait()

		return loadBalancers, errors
	},
	SetLabelsBatch: func(s state.State, loadBalancers []*hcloud.LoadBalancer, labels map[string]string) []error {
		errors := make([]error, len(loadBalancers))

		var wg sync.WaitGroup
		for i, loadBalancer := range loadBalancers {
			if loadBalancer == nil {
				continue
			}

			wg.Add(1)
			go func(idx int, lb *hcloud.LoadBalancer) {
				defer wg.Done()
				opts := hcloud.LoadBalancerUpdateOpts{
					Labels: labels,
				}
				_, _, err := s.Client().LoadBalancer().Update(s, lb, opts)
				errors[idx] = err
			}(i, loadBalancer)
		}
		wg.Wait()

		return errors
	},
}

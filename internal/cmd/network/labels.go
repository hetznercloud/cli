package network

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var LabelCmds = base.LabelCmds[*hcloud.Network]{
	ResourceNameSingular:   "Network",
	ShortDescriptionAdd:    "Add a label to a Network",
	ShortDescriptionRemove: "Remove a label from a Network",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.Network().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.Network().LabelKeys },
	Fetch: func(s state.State, idOrName string) (*hcloud.Network, error) {
		network, _, err := s.Client().Network().Get(s, idOrName)
		if err != nil {
			return nil, err
		}
		if network == nil {
			return nil, fmt.Errorf("Network not found: %s", idOrName)
		}
		return network, nil
	},
	SetLabels: func(s state.State, network *hcloud.Network, labels map[string]string) error {
		opts := hcloud.NetworkUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Network().Update(s, network, opts)
		return err
	},
	GetLabels: func(network *hcloud.Network) map[string]string {
		return network.Labels
	},
	GetIDOrName: func(network *hcloud.Network) string {
		return strconv.FormatInt(network.ID, 10)
	},
	FetchBatch: func(s state.State, idOrNames []string) ([]*hcloud.Network, []error) {
		networks := make([]*hcloud.Network, len(idOrNames))
		errors := make([]error, len(idOrNames))

		var wg sync.WaitGroup
		for i, idOrName := range idOrNames {
			wg.Add(1)
			go func(idx int, id string) {
				defer wg.Done()
				network, _, err := s.Client().Network().Get(s, id)
				if err != nil {
					errors[idx] = err
					return
				}
				if network == nil {
					errors[idx] = fmt.Errorf("Network not found: %s", id)
					return
				}
				networks[idx] = network
			}(i, idOrName)
		}
		wg.Wait()

		return networks, errors
	},
	SetLabelsBatch: func(s state.State, networks []*hcloud.Network, labels map[string]string) []error {
		errors := make([]error, len(networks))

		var wg sync.WaitGroup
		for i, network := range networks {
			if network == nil {
				continue
			}

			wg.Add(1)
			go func(idx int, net *hcloud.Network) {
				defer wg.Done()
				opts := hcloud.NetworkUpdateOpts{
					Labels: labels,
				}
				_, _, err := s.Client().Network().Update(s, net, opts)
				errors[idx] = err
			}(i, network)
		}
		wg.Wait()

		return errors
	},
}

package network

import (
	"fmt"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
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
}

package network

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds{
	ResourceNameSingular:   "Network",
	ShortDescriptionAdd:    "Add a label to a Network",
	ShortDescriptionRemove: "Remove a label from a Network",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.Network().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.Network().LabelKeys },
	FetchLabels: func(s state.State, idOrName string) (map[string]string, int64, error) {
		network, _, err := s.Client().Network().Get(s, idOrName)
		if err != nil {
			return nil, 0, err
		}
		if network == nil {
			return nil, 0, fmt.Errorf("network not found: %s", idOrName)
		}
		return network.Labels, network.ID, nil
	},
	SetLabels: func(s state.State, id int64, labels map[string]string) error {
		opts := hcloud.NetworkUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Network().Update(s, &hcloud.Network{ID: id}, opts)
		return err
	},
}

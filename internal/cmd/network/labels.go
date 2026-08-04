package network

import (
	"context"
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
	NameSuggestions:        func(c hcapi2.Client) hcapi2.CompletionFunc { return c.Network().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) hcapi2.LabelCompletionFunc { return c.Network().LabelKeys },
	Fetch: func(ctx context.Context, s state.State, idOrName string) (*hcloud.Network, error) {
		network, _, err := s.Client().Network().Get(ctx, idOrName)
		if err != nil {
			return nil, err
		}
		if network == nil {
			return nil, fmt.Errorf("Network not found: %s", idOrName)
		}
		return network, nil
	},
	SetLabels: func(ctx context.Context, s state.State, network *hcloud.Network, labels map[string]string) error {
		opts := hcloud.NetworkUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Network().Update(ctx, network, opts)
		return err
	},
	GetLabels: func(network *hcloud.Network) map[string]string {
		return network.Labels
	},
	GetIDOrName: func(network *hcloud.Network) string {
		return strconv.FormatInt(network.ID, 10)
	},
}

package network

import (
	"context"
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds{
	ResourceNameSingular:   "Network",
	ShortDescriptionAdd:    "Add a label to a Network",
	ShortDescriptionRemove: "Remove a label from a Network",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.Network().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.Network().LabelKeys },
	FetchLabels: func(ctx context.Context, client hcapi2.Client, idOrName string) (map[string]string, int64, error) {
		network, _, err := client.Network().Get(ctx, idOrName)
		if err != nil {
			return nil, 0, err
		}
		if network == nil {
			return nil, 0, fmt.Errorf("network not found: %s", idOrName)
		}
		return network.Labels, network.ID, nil
	},
	SetLabels: func(ctx context.Context, client hcapi2.Client, id int64, labels map[string]string) error {
		opts := hcloud.NetworkUpdateOpts{
			Labels: labels,
		}
		_, _, err := client.Network().Update(ctx, &hcloud.Network{ID: id}, opts)
		return err
	},
}

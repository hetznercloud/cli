package firewall

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds{
	ResourceNameSingular:   "firewall",
	ShortDescriptionAdd:    "Add a label to an firewall",
	ShortDescriptionRemove: "Remove a label from an firewall",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.Firewall().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.Firewall().LabelKeys },
	FetchLabels: func(s state.State, idOrName string) (map[string]string, int64, error) {
		firewall, _, err := s.Client().Firewall().Get(s, idOrName)
		if err != nil {
			return nil, 0, err
		}
		if firewall == nil {
			return nil, 0, fmt.Errorf("firewall not found: %s", idOrName)
		}
		return firewall.Labels, firewall.ID, nil
	},
	SetLabels: func(s state.State, id int64, labels map[string]string) error {
		opts := hcloud.FirewallUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Firewall().Update(s, &hcloud.Firewall{ID: id}, opts)
		return err
	},
}

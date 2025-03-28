package firewall

import (
	"fmt"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds[*hcloud.Firewall]{
	ResourceNameSingular:   "Firewall",
	ShortDescriptionAdd:    "Add a label to a Firewall",
	ShortDescriptionRemove: "Remove a label from a Firewall",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.Firewall().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.Firewall().LabelKeys },
	Fetch: func(s state.State, idOrName string) (*hcloud.Firewall, error) {
		firewall, _, err := s.Client().Firewall().Get(s, idOrName)
		if err != nil {
			return nil, err
		}
		if firewall == nil {
			return nil, fmt.Errorf("Firewall not found: %s", idOrName)
		}
		return firewall, nil
	},
	SetLabels: func(s state.State, firewall *hcloud.Firewall, labels map[string]string) error {
		opts := hcloud.FirewallUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Firewall().Update(s, firewall, opts)
		return err
	},
	GetLabels: func(firewall *hcloud.Firewall) map[string]string {
		return firewall.Labels
	},
	GetIDOrName: func(firewall *hcloud.Firewall) string {
		return strconv.FormatInt(firewall.ID, 10)
	},
}

package firewall

import (
	"fmt"
	"strconv"

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
	Fetch: func(s state.State, idOrName string) (any, error) {
		firewall, _, err := s.Client().Firewall().Get(s, idOrName)
		if err != nil {
			return nil, err
		}
		if firewall == nil {
			return nil, fmt.Errorf("firewall not found: %s", idOrName)
		}
		return firewall, nil
	},
	SetLabels: func(s state.State, resource any, labels map[string]string) error {
		firewall := resource.(*hcloud.Firewall)
		opts := hcloud.FirewallUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Firewall().Update(s, firewall, opts)
		return err
	},
	GetLabels: func(resource any) map[string]string {
		firewall := resource.(*hcloud.Firewall)
		return firewall.Labels
	},
	GetIDOrName: func(resource any) string {
		firewall := resource.(*hcloud.Firewall)
		return strconv.FormatInt(firewall.ID, 10)
	},
}

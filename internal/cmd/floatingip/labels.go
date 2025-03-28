package floatingip

import (
	"fmt"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds[*hcloud.FloatingIP]{
	ResourceNameSingular:   "Floating IP",
	ShortDescriptionAdd:    "Add a label to a Floating IP",
	ShortDescriptionRemove: "Remove a label from a Floating IP",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.FloatingIP().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.FloatingIP().LabelKeys },
	Fetch: func(s state.State, idOrName string) (*hcloud.FloatingIP, error) {
		floatingIP, _, err := s.Client().FloatingIP().Get(s, idOrName)
		if err != nil {
			return nil, err
		}
		if floatingIP == nil {
			return nil, fmt.Errorf("Floating IP not found: %s", idOrName)
		}
		return floatingIP, nil
	},
	SetLabels: func(s state.State, floatingIP *hcloud.FloatingIP, labels map[string]string) error {
		opts := hcloud.FloatingIPUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().FloatingIP().Update(s, floatingIP, opts)
		return err
	},
	GetLabels: func(floatingIP *hcloud.FloatingIP) map[string]string {
		return floatingIP.Labels
	},
	GetIDOrName: func(floatingIP *hcloud.FloatingIP) string {
		return strconv.FormatInt(floatingIP.ID, 10)
	},
}

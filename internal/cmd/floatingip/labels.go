package floatingip

import (
	"fmt"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds{
	ResourceNameSingular:   "Floating IP",
	ShortDescriptionAdd:    "Add a label to an Floating IP",
	ShortDescriptionRemove: "Remove a label from an Floating IP",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.FloatingIP().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.FloatingIP().LabelKeys },
	Fetch: func(s state.State, idOrName string) (any, error) {
		floatingIP, _, err := s.Client().FloatingIP().Get(s, idOrName)
		if err != nil {
			return nil, err
		}
		if floatingIP == nil {
			return nil, fmt.Errorf("floating IP not found: %s", idOrName)
		}
		return floatingIP, nil
	},
	SetLabels: func(s state.State, resource any, labels map[string]string) error {
		floatingIP := resource.(*hcloud.FloatingIP)
		opts := hcloud.FloatingIPUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().FloatingIP().Update(s, floatingIP, opts)
		return err
	},
	GetLabels: func(resource any) map[string]string {
		floatingIP := resource.(*hcloud.FloatingIP)
		return floatingIP.Labels
	},
	GetIDOrName: func(resource any) string {
		floatingIP := resource.(*hcloud.FloatingIP)
		return strconv.FormatInt(floatingIP.ID, 10)
	},
}

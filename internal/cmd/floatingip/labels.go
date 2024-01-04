package floatingip

import (
	"fmt"

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
	FetchLabels: func(s state.State, idOrName string) (map[string]string, int64, error) {
		floatingIP, _, err := s.Client().FloatingIP().Get(s, idOrName)
		if err != nil {
			return nil, 0, err
		}
		if floatingIP == nil {
			return nil, 0, fmt.Errorf("floating IP not found: %s", idOrName)
		}
		return floatingIP.Labels, floatingIP.ID, nil
	},
	SetLabels: func(s state.State, id int64, labels map[string]string) error {
		opts := hcloud.FloatingIPUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().FloatingIP().Update(s, &hcloud.FloatingIP{ID: id}, opts)
		return err
	},
}

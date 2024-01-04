package primaryip

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds{
	ResourceNameSingular:   "primary-ip",
	ShortDescriptionAdd:    "Add a label to a Primary IP",
	ShortDescriptionRemove: "Remove a label from a Primary IP",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.PrimaryIP().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.PrimaryIP().LabelKeys },
	FetchLabels: func(s state.State, idOrName string) (map[string]string, int64, error) {
		primaryIP, _, err := s.Client().PrimaryIP().Get(s, idOrName)
		if err != nil {
			return nil, 0, err
		}
		if primaryIP == nil {
			return nil, 0, fmt.Errorf("primaryIP not found: %s", idOrName)
		}
		return primaryIP.Labels, primaryIP.ID, nil
	},
	SetLabels: func(s state.State, id int64, labels map[string]string) error {
		opts := hcloud.PrimaryIPUpdateOpts{
			Labels: &labels,
		}
		_, _, err := s.Client().PrimaryIP().Update(s, &hcloud.PrimaryIP{ID: id}, opts)
		return err
	},
}

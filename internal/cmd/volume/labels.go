package volume

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds{
	ResourceNameSingular:   "Volume",
	ShortDescriptionAdd:    "Add a label to a Volume",
	ShortDescriptionRemove: "Remove a label from a Volume",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.Volume().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.Volume().LabelKeys },
	FetchLabels: func(s state.State, idOrName string) (map[string]string, int64, error) {
		volume, _, err := s.Client().Volume().Get(s, idOrName)
		if err != nil {
			return nil, 0, err
		}
		if volume == nil {
			return nil, 0, fmt.Errorf("volume not found: %s", idOrName)
		}
		return volume.Labels, volume.ID, nil
	},
	SetLabels: func(s state.State, id int64, labels map[string]string) error {
		opts := hcloud.VolumeUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Volume().Update(s, &hcloud.Volume{ID: id}, opts)
		return err
	},
}

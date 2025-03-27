package volume

import (
	"fmt"
	"strconv"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds[*hcloud.Volume]{
	ResourceNameSingular:   "Volume",
	ShortDescriptionAdd:    "Add a label to a Volume",
	ShortDescriptionRemove: "Remove a label from a Volume",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.Volume().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.Volume().LabelKeys },
	Fetch: func(s state.State, idOrName string) (*hcloud.Volume, error) {
		volume, _, err := s.Client().Volume().Get(s, idOrName)
		if err != nil {
			return nil, err
		}
		if volume == nil {
			return nil, fmt.Errorf("Volume not found: %s", idOrName)
		}
		return volume, nil
	},
	SetLabels: func(s state.State, volume *hcloud.Volume, labels map[string]string) error {
		opts := hcloud.VolumeUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().Volume().Update(s, volume, opts)
		return err
	},
	GetLabels: func(volume *hcloud.Volume) map[string]string {
		return volume.Labels
	},
	GetIDOrName: func(volume *hcloud.Volume) string {
		return strconv.FormatInt(volume.ID, 10)
	},
}

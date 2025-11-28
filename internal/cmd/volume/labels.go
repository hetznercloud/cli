package volume

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
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
	FetchBatch: func(s state.State, idOrNames []string) ([]*hcloud.Volume, []error) {
		volumes := make([]*hcloud.Volume, len(idOrNames))
		errors := make([]error, len(idOrNames))

		var wg sync.WaitGroup
		for i, idOrName := range idOrNames {
			wg.Add(1)
			go func(idx int, id string) {
				defer wg.Done()
				volume, _, err := s.Client().Volume().Get(s, id)
				if err != nil {
					errors[idx] = err
					return
				}
				if volume == nil {
					errors[idx] = fmt.Errorf("Volume not found: %s", id)
					return
				}
				volumes[idx] = volume
			}(i, idOrName)
		}
		wg.Wait()

		return volumes, errors
	},
	SetLabelsBatch: func(s state.State, volumes []*hcloud.Volume, labels map[string]string) []error {
		errors := make([]error, len(volumes))

		var wg sync.WaitGroup
		for i, volume := range volumes {
			if volume == nil {
				continue
			}

			wg.Add(1)
			go func(idx int, vol *hcloud.Volume) {
				defer wg.Done()
				opts := hcloud.VolumeUpdateOpts{
					Labels: labels,
				}
				_, _, err := s.Client().Volume().Update(s, vol, opts)
				errors[idx] = err
			}(i, volume)
		}
		wg.Wait()

		return errors
	},
}

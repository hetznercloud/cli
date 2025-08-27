package storagebox

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds[*hcloud.StorageBox]{
	ResourceNameSingular:   "Storage Box",
	ShortDescriptionAdd:    "Add a label to a Storage Box",
	ShortDescriptionRemove: "Remove a label from a Storage Box",
	NameSuggestions:        func(c hcapi2.Client) func() []string { return c.StorageBox().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) func(idOrName string) []string { return c.StorageBox().LabelKeys },
	Fetch: func(s state.State, idOrName string) (*hcloud.StorageBox, error) {
		storageBox, _, err := s.Client().StorageBox().Get(s, idOrName)
		if err != nil {
			return nil, err
		}
		if storageBox == nil {
			return nil, fmt.Errorf("Storage Box not found: %s", idOrName)
		}
		return storageBox, nil
	},
	SetLabels: func(s state.State, storageBox *hcloud.StorageBox, labels map[string]string) error {
		opts := hcloud.StorageBoxUpdateOpts{
			Labels: labels,
			// TODO: why do we need to re-set the name?
			Name: storageBox.Name,
		}
		_, _, err := s.Client().StorageBox().Update(s, storageBox, opts)
		return err
	},
	GetLabels: func(storageBox *hcloud.StorageBox) map[string]string {
		return storageBox.Labels
	},
	GetIDOrName: func(storageBox *hcloud.StorageBox) string {
		return storageBox.Name
	},
}

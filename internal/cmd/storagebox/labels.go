package storagebox

import (
	"context"
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
	NameSuggestions:        func(c hcapi2.Client) hcapi2.CompletionFunc { return c.StorageBox().Names },
	LabelKeySuggestions:    func(c hcapi2.Client) hcapi2.LabelCompletionFunc { return c.StorageBox().LabelKeys },
	Fetch: func(ctx context.Context, s state.State, idOrName string) (*hcloud.StorageBox, error) {
		storageBox, _, err := s.Client().StorageBox().Get(ctx, idOrName)
		if err != nil {
			return nil, err
		}
		if storageBox == nil {
			return nil, fmt.Errorf("Storage Box not found: %s", idOrName)
		}
		return storageBox, nil
	},
	SetLabels: func(ctx context.Context, s state.State, storageBox *hcloud.StorageBox, labels map[string]string) error {
		opts := hcloud.StorageBoxUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().StorageBox().Update(ctx, storageBox, opts)
		return err
	},
	GetLabels: func(storageBox *hcloud.StorageBox) map[string]string {
		return storageBox.Labels
	},
	GetIDOrName: func(storageBox *hcloud.StorageBox) string {
		return storageBox.Name
	},
}

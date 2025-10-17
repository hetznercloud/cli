package snapshot

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var LabelCmds = base.LabelCmds[*hcloud.StorageBoxSnapshot]{
	ResourceNameSingular:   "Storage Box Snapshot",
	ShortDescriptionAdd:    "Add a label to a Storage Box Snapshot",
	ShortDescriptionRemove: "Remove a label from a Storage Box Snapshot",

	PositionalArgumentOverride: []string{"storage-box", "snapshot"},
	ValidArgsFunction: func(client hcapi2.Client) []cobra.CompletionFunc {
		return []cobra.CompletionFunc{
			cmpl.SuggestCandidatesF(client.StorageBox().Names),
			SuggestSnapshots(client),
			cmpl.SuggestCandidatesCtx(func(_ *cobra.Command, args []string) []string {
				if len(args) < 2 {
					return nil
				}
				return client.StorageBox().SnapshotLabelKeys(args[0], args[1])
			}),
		}
	},

	FetchWithArgs: func(s state.State, args []string) (*hcloud.StorageBoxSnapshot, error) {
		storageBoxIDOrName := args[0]

		storageBox, _, err := s.Client().StorageBox().Get(s, storageBoxIDOrName)
		if err != nil {
			return nil, err
		}
		if storageBox == nil {
			return nil, fmt.Errorf("Storage Box not found: %s", storageBoxIDOrName)
		}

		snapshot, _, err := s.Client().StorageBox().GetSnapshot(s, storageBox, args[1])
		if err != nil {
			return nil, err
		}
		return snapshot, nil
	},

	SetLabels: func(s state.State, snapshot *hcloud.StorageBoxSnapshot, labels map[string]string) error {
		opts := hcloud.StorageBoxSnapshotUpdateOpts{
			Labels: labels,
		}
		_, _, err := s.Client().StorageBox().UpdateSnapshot(s, snapshot, opts)
		return err
	},

	GetLabels: func(snapshot *hcloud.StorageBoxSnapshot) map[string]string {
		return snapshot.Labels
	},

	GetIDOrName: func(snapshot *hcloud.StorageBoxSnapshot) string {
		return snapshot.Name
	},
}

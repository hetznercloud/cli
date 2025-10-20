package snapshot

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.StorageBoxSnapshot]{
	ResourceNameSingular:       "Storage Box Snapshot",
	ShortDescription:           "Describe a Storage Box Snapshot",
	PositionalArgumentOverride: []string{"storage-box", "snapshot"},
	ValidArgsFunction: func(client hcapi2.Client) []cobra.CompletionFunc {
		return []cobra.CompletionFunc{
			cmpl.SuggestCandidatesF(client.StorageBox().Names),
			SuggestSnapshots(client),
		}
	},
	FetchWithArgs: func(s state.State, _ *cobra.Command, args []string) (*hcloud.StorageBoxSnapshot, any, error) {
		storageBoxIDOrName := args[0]

		storageBox, _, err := s.Client().StorageBox().Get(s, storageBoxIDOrName)
		if err != nil {
			return nil, nil, err
		}
		if storageBox == nil {
			return nil, nil, fmt.Errorf("Storage Box not found: %s", storageBoxIDOrName)
		}

		snapshot, _, err := s.Client().StorageBox().GetSnapshot(s, storageBox, args[1])
		if err != nil {
			return nil, nil, err
		}
		return snapshot, hcloud.SchemaFromStorageBoxSnapshot(snapshot), nil
	},
	PrintText: func(_ state.State, cmd *cobra.Command, snapshot *hcloud.StorageBoxSnapshot) error {
		cmd.Printf("ID:\t\t\t%d\n", snapshot.ID)
		cmd.Printf("Name:\t\t\t%s\n", snapshot.Name)
		cmd.Printf("Description:\t\t%s\n", snapshot.Description)
		cmd.Printf("Created:\t\t%s (%s)\n", util.Datetime(snapshot.Created), humanize.Time(snapshot.Created))
		cmd.Printf("Is automatic:\t\t%s\n", util.YesNo(snapshot.IsAutomatic))

		cmd.Println("Stats:")
		cmd.Printf("  Size:\t\t\t%s\n", humanize.IBytes(snapshot.Stats.Size))
		cmd.Printf("  Filesystem Size:\t%s\n", humanize.IBytes(snapshot.Stats.SizeFilesystem))

		cmd.Println("Labels:")
		if len(snapshot.Labels) == 0 {
			cmd.Println("  No labels")
		} else {
			for key, value := range util.IterateInOrder(snapshot.Labels) {
				cmd.Printf("  %s: %s\n", key, value)
			}
		}

		cmd.Println("Storage Box:")
		cmd.Printf("  ID:\t\t\t%d\n", snapshot.StorageBox.ID)
		return nil
	},
}

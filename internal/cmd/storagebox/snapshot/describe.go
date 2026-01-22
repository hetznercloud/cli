package snapshot

import (
	"fmt"
	"io"

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
	PrintText: func(_ state.State, _ *cobra.Command, out io.Writer, snapshot *hcloud.StorageBoxSnapshot) error {
		fmt.Fprintf(out, "ID:\t%d\n", snapshot.ID)
		fmt.Fprintf(out, "Name:\t%s\n", snapshot.Name)
		fmt.Fprintf(out, "Description:\t%s\n", snapshot.Description)
		fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(snapshot.Created), humanize.Time(snapshot.Created))
		fmt.Fprintf(out, "Is automatic:\t%s\n", util.YesNo(snapshot.IsAutomatic))

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Stats:\n")
		fmt.Fprintf(out, "  Size:\t%s\n", humanize.IBytes(snapshot.Stats.Size))
		fmt.Fprintf(out, "  Filesystem Size:\t%s\n", humanize.IBytes(snapshot.Stats.SizeFilesystem))

		fmt.Fprintln(out)
		util.DescribeLabels(out, snapshot.Labels, "")

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Storage Box:\n")
		fmt.Fprintf(out, "  ID:\t%d\n", snapshot.StorageBox.ID)

		return nil
	},
}

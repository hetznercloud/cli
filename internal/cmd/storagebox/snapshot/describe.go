package snapshot

import (
	"fmt"
	"io"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
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
		_, _ = fmt.Fprint(out, DescribeSnapshot(snapshot))
		return nil
	},
	Experimental: experimental.StorageBoxes,
}

func DescribeSnapshot(snapshot *hcloud.StorageBoxSnapshot) string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, "ID:\t%d\n", snapshot.ID)
	_, _ = fmt.Fprintf(&sb, "Name:\t%s\n", snapshot.Name)
	_, _ = fmt.Fprintf(&sb, "Description:\t%s\n", snapshot.Description)
	_, _ = fmt.Fprintf(&sb, "Created:\t%s (%s)\n", util.Datetime(snapshot.Created), humanize.Time(snapshot.Created))
	_, _ = fmt.Fprintf(&sb, "Is automatic:\t%s\n", util.YesNo(snapshot.IsAutomatic))

	_, _ = fmt.Fprintf(&sb, "Stats:\n")
	_, _ = fmt.Fprintf(&sb, "  Size:\t%s\n", humanize.IBytes(snapshot.Stats.Size))
	_, _ = fmt.Fprintf(&sb, "  Filesystem Size:\t%s\n", humanize.IBytes(snapshot.Stats.SizeFilesystem))

	util.DescribeLabels(&sb, snapshot.Labels, "")

	_, _ = fmt.Fprintf(&sb, "Storage Box:\n")
	_, _ = fmt.Fprintf(&sb, "  ID:\t%d\n", snapshot.StorageBox.ID)
	return sb.String()
}

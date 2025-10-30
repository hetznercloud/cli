package storagebox

import (
	"fmt"
	"io"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
	"github.com/hetznercloud/cli/internal/cmd/location"
	"github.com/hetznercloud/cli/internal/cmd/storageboxtype"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.StorageBox]{
	ResourceNameSingular: "Storage Box",
	ShortDescription:     "Describe a Storage Box",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.StorageBox().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.StorageBox, any, error) {
		storageBox, _, err := s.Client().StorageBox().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return storageBox, hcloud.SchemaFromStorageBox(storageBox), nil
	},
	PrintText: func(s state.State, _ *cobra.Command, out io.Writer, storageBox *hcloud.StorageBox) error {

		_, _ = fmt.Fprintf(out, "ID:\t%d\n", storageBox.ID)
		_, _ = fmt.Fprintf(out, "Name:\t%s\n", storageBox.Name)
		_, _ = fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(storageBox.Created), humanize.Time(storageBox.Created))
		_, _ = fmt.Fprintf(out, "Status:\t%s\n", storageBox.Status)
		_, _ = fmt.Fprintf(out, "Username:\t%s\n", storageBox.Username)
		_, _ = fmt.Fprintf(out, "Server:\t%s\n", storageBox.Server)
		_, _ = fmt.Fprintf(out, "System:\t%s\n", storageBox.System)

		snapshotPlan := storageBox.SnapshotPlan
		if snapshotPlan == nil {
			_, _ = fmt.Fprintf(out, "Snapshot Plan:\tNo snapshot plan active\n")
		} else {
			_, _ = fmt.Fprintf(out, "Snapshot Plan:\t\n")
			_, _ = fmt.Fprintf(out, "  Max Snapshots:\t%d\n", snapshotPlan.MaxSnapshots)
			_, _ = fmt.Fprintf(out, "  Minute:\t%d\n", snapshotPlan.Minute)
			_, _ = fmt.Fprintf(out, "  Hour:\t%d\n", snapshotPlan.Hour)

			if snapshotPlan.DayOfWeek != nil {
				_, _ = fmt.Fprintf(out, "  Day of Week:\t%s\n", *snapshotPlan.DayOfWeek)
			}
			if snapshotPlan.DayOfMonth != nil {
				_, _ = fmt.Fprintf(out, "  Day of Month:\t%d\n", *snapshotPlan.DayOfMonth)
			}
		}

		protection := storageBox.Protection
		_, _ = fmt.Fprintf(out, "Protection:\t\n")
		_, _ = fmt.Fprintf(out, "  Delete:\t%t\n", protection.Delete)

		stats := storageBox.Stats
		_, _ = fmt.Fprintf(out, "Stats:\t\n")
		_, _ = fmt.Fprintf(out, "  Size:\t%s\n", humanize.IBytes(stats.Size))
		_, _ = fmt.Fprintf(out, "  Size Data:\t%s\n", humanize.IBytes(stats.SizeData))
		_, _ = fmt.Fprintf(out, "  Size Snapshots:\t%s\n", humanize.IBytes(stats.SizeSnapshots))

		util.DescribeLabels(out, storageBox.Labels, "")

		accessSettings := storageBox.AccessSettings
		_, _ = fmt.Fprintf(out, "Access Settings:\t\n")
		_, _ = fmt.Fprintf(out, "  Reachable Externally:\t%t\n", accessSettings.ReachableExternally)
		_, _ = fmt.Fprintf(out, "  Samba Enabled:\t%t\n", accessSettings.SambaEnabled)
		_, _ = fmt.Fprintf(out, "  SSH Enabled:\t%t\n", accessSettings.SSHEnabled)
		_, _ = fmt.Fprintf(out, "  WebDAV Enabled:\t%t\n", accessSettings.WebDAVEnabled)
		_, _ = fmt.Fprintf(out, "  ZFS Enabled:\t%t\n", accessSettings.ZFSEnabled)

		typeDescription, _ := storageboxtype.DescribeStorageBoxType(s, storageBox.StorageBoxType, true)
		_, _ = fmt.Fprintf(out, "Storage Box Type:\t\n")
		_, _ = fmt.Fprintf(out, "%s", util.PrefixLines(typeDescription, "  "))

		_, _ = fmt.Fprintf(out, "Location:\t\n")
		_, _ = fmt.Fprintf(out, "%s", util.PrefixLines(location.DescribeLocation(storageBox.Location), "  "))

		return nil
	},
	Experimental: experimental.StorageBoxes,
}

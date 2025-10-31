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

		fmt.Fprintf(out, "ID:\t%d\n", storageBox.ID)
		fmt.Fprintf(out, "Name:\t%s\n", storageBox.Name)
		fmt.Fprintf(out, "Created:\t%s (%s)\n", util.Datetime(storageBox.Created), humanize.Time(storageBox.Created))
		fmt.Fprintf(out, "Status:\t%s\n", storageBox.Status)
		fmt.Fprintf(out, "Username:\t%s\n", storageBox.Username)
		fmt.Fprintf(out, "Server:\t%s\n", storageBox.Server)
		fmt.Fprintf(out, "System:\t%s\n", storageBox.System)

		snapshotPlan := storageBox.SnapshotPlan
		fmt.Fprintln(out)
		fmt.Fprintf(out, "Snapshot Plan:\n")
		if snapshotPlan == nil {
			fmt.Fprintf(out, "  No snapshot plan active\n")
		} else {
			fmt.Fprintf(out, "  Max Snapshots:\t%d\n", snapshotPlan.MaxSnapshots)
			fmt.Fprintf(out, "  Minute:\t%d\n", snapshotPlan.Minute)
			fmt.Fprintf(out, "  Hour:\t%d\n", snapshotPlan.Hour)

			if snapshotPlan.DayOfWeek != nil {
				fmt.Fprintf(out, "  Day of Week:\t%s\n", *snapshotPlan.DayOfWeek)
			}
			if snapshotPlan.DayOfMonth != nil {
				fmt.Fprintf(out, "  Day of Month:\t%d\n", *snapshotPlan.DayOfMonth)
			}
		}

		protection := storageBox.Protection
		fmt.Fprintln(out)
		fmt.Fprintf(out, "Protection:\n")
		fmt.Fprintf(out, "  Delete:\t%t\n", protection.Delete)

		stats := storageBox.Stats
		fmt.Fprintln(out)
		fmt.Fprintf(out, "Stats:\n")
		fmt.Fprintf(out, "  Size:\t%s\n", humanize.IBytes(stats.Size))
		fmt.Fprintf(out, "  Size Data:\t%s\n", humanize.IBytes(stats.SizeData))
		fmt.Fprintf(out, "  Size Snapshots:\t%s\n", humanize.IBytes(stats.SizeSnapshots))

		fmt.Fprintln(out)
		util.DescribeLabels(out, storageBox.Labels, "")

		accessSettings := storageBox.AccessSettings
		fmt.Fprintln(out)
		fmt.Fprintf(out, "Access Settings:\n")
		fmt.Fprintf(out, "  Reachable Externally:\t%t\n", accessSettings.ReachableExternally)
		fmt.Fprintf(out, "  Samba Enabled:\t%t\n", accessSettings.SambaEnabled)
		fmt.Fprintf(out, "  SSH Enabled:\t%t\n", accessSettings.SSHEnabled)
		fmt.Fprintf(out, "  WebDAV Enabled:\t%t\n", accessSettings.WebDAVEnabled)
		fmt.Fprintf(out, "  ZFS Enabled:\t%t\n", accessSettings.ZFSEnabled)

		typeDescription, _ := storageboxtype.DescribeStorageBoxType(s, storageBox.StorageBoxType, true)
		fmt.Fprintln(out)
		fmt.Fprintf(out, "Storage Box Type:\n")
		fmt.Fprintf(out, "%s", util.PrefixLines(typeDescription, "  "))

		fmt.Fprintln(out)
		fmt.Fprintf(out, "Location:\t\n")
		fmt.Fprintf(out, "%s", util.PrefixLines(location.DescribeLocation(storageBox.Location), "  "))

		return nil
	},
	Experimental: experimental.StorageBoxes,
}

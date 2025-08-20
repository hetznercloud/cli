package storagebox

import (
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
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
	PrintText: func(_ state.State, cmd *cobra.Command, storageBox *hcloud.StorageBox) error {
		cmd.Printf("ID:\t\t\t\t%d\n", storageBox.ID)
		cmd.Printf("Name:\t\t\t\t%s\n", storageBox.Name)
		cmd.Printf("Created:\t\t\t%s (%s)\n", util.Datetime(storageBox.Created), humanize.Time(storageBox.Created))
		cmd.Printf("Status:\t\t\t\t%s\n", storageBox.Status)
		cmd.Printf("Username:\t\t\t%s\n", util.OptionalString(storageBox.Username, "-"))
		cmd.Printf("Server:\t\t\t\t%s\n", util.OptionalString(storageBox.Server, "-"))
		cmd.Printf("System:\t\t\t\t%s\n", util.OptionalString(storageBox.System, "-"))

		snapshotPlan := storageBox.SnapshotPlan
		cmd.Println("Snapshot Plan:")
		if snapshotPlan == nil {
			cmd.Println("  No snapshot plan available")
		} else {
			cmd.Printf("  Max Snapshots:\t\t%d\n", snapshotPlan.MaxSnapshots)
			if snapshotPlan.Minute != nil {
				cmd.Printf("  Minute:\t\t\t%d\n", *snapshotPlan.Minute)
			}
			if snapshotPlan.Hour != nil {
				cmd.Printf("  Hour:\t\t\t\t%d\n", *snapshotPlan.Hour)
			}
			if snapshotPlan.DayOfWeek != nil {
				cmd.Printf("  Day of Week:\t\t\t%d\n", *snapshotPlan.DayOfWeek)
			}
			if snapshotPlan.DayOfMonth != nil {
				cmd.Printf("  Day of Month:\t\t\t%d\n", *snapshotPlan.DayOfMonth)
			}
		}

		protection := storageBox.Protection
		cmd.Println("Protection:")
		cmd.Printf("  Delete:\t\t\t%t\n", protection.Delete)

		stats := storageBox.Stats
		cmd.Println("Stats:")
		if stats == nil {
			cmd.Println("  No stats available")
		} else {
			cmd.Printf("  Size:\t\t\t\t%s\n", humanize.IBytes(stats.Size))
			cmd.Printf("  Size Data:\t\t\t%s\n", humanize.IBytes(stats.SizeData))
			cmd.Printf("  Size Snapshots:\t\t%s\n", humanize.IBytes(stats.SizeSnapshots))
		}

		cmd.Print("Labels:\n")
		if len(storageBox.Labels) == 0 {
			cmd.Print("  No labels\n")
		} else {
			for key, value := range util.IterateInOrder(storageBox.Labels) {
				cmd.Printf("  %s: %s\n", key, value)
			}
		}

		accessSettings := storageBox.AccessSettings
		cmd.Println("Access Settings:")
		cmd.Printf("  Reachable Externally:\t\t%t\n", accessSettings.ReachableExternally)
		cmd.Printf("  Samba Enabled:\t\t%t\n", accessSettings.SambaEnabled)
		cmd.Printf("  SSH Enabled:\t\t\t%t\n", accessSettings.SSHEnabled)
		cmd.Printf("  WebDAV Enabled:\t\t%t\n", accessSettings.WebDAVEnabled)
		cmd.Printf("  ZFS Enabled:\t\t\t%t\n", accessSettings.ZFSEnabled)

		storageBoxType := storageBox.StorageBoxType
		cmd.Println("Storage Box Type:")
		cmd.Printf("  ID:\t\t\t\t%d\n", storageBoxType.ID)
		cmd.Printf("  Name:\t\t\t\t%s\n", storageBoxType.Name)
		cmd.Printf("  Description:\t\t\t%s\n", storageBoxType.Description)
		cmd.Printf("  Size:\t\t\t\t%s\n", humanize.IBytes(uint64(storageBoxType.Size)))
		cmd.Printf("  Snapshot Limit:\t\t%d\n", storageBoxType.SnapshotLimit)
		cmd.Printf("  Automatic Snapshot Limit:\t%d\n", storageBoxType.AutomaticSnapshotLimit)
		cmd.Printf("  Subaccounts Limit:\t\t%d\n", storageBoxType.SubaccountsLimit)

		location := storageBox.Location
		cmd.Println("Location:")
		cmd.Printf("  ID:\t\t\t\t%d\n", location.ID)
		cmd.Printf("  Name:\t\t\t\t%s\n", location.Name)
		cmd.Printf("  Description:\t\t\t%s\n", location.Description)
		cmd.Printf("  Network Zone:\t\t\t%s\n", location.NetworkZone)
		cmd.Printf("  Country:\t\t\t%s\n", location.Country)
		cmd.Printf("  City:\t\t\t\t%s\n", location.City)
		cmd.Printf("  Latitude:\t\t\t%f\n", location.Latitude)
		cmd.Printf("  Longitude:\t\t\t%f\n", location.Longitude)

		return nil
	},
}

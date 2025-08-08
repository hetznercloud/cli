package storageboxtype

import (
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var DescribeCmd = base.DescribeCmd[*hcloud.StorageBoxType]{
	ResourceNameSingular: "Storage Box Type",
	ShortDescription:     "Describe a Storage Box Type",
	NameSuggestions:      func(c hcapi2.Client) func() []string { return c.StorageBoxType().Names },
	Fetch: func(s state.State, _ *cobra.Command, idOrName string) (*hcloud.StorageBoxType, any, error) {
		st, _, err := s.Client().StorageBoxType().Get(s, idOrName)
		if err != nil {
			return nil, nil, err
		}
		return st, hcloud.SchemaFromStorageBoxType(st), nil
	},
	PrintText: func(s state.State, cmd *cobra.Command, storageBoxType *hcloud.StorageBoxType) error {
		cmd.Printf("ID:\t\t\t\t%d\n", storageBoxType.ID)
		cmd.Printf("Name:\t\t\t\t%s\n", storageBoxType.Name)
		cmd.Printf("Description:\t\t\t%s\n", storageBoxType.Description)
		cmd.Printf("Size:\t\t\t\t%s\n", humanize.IBytes(uint64(storageBoxType.Size)))
		cmd.Printf("Snapshot Limit:\t\t\t%d\n", storageBoxType.SnapshotLimit)
		cmd.Printf("Automatic Snapshot Limit:\t%d\n", storageBoxType.AutomaticSnapshotLimit)
		cmd.Printf("Subaccounts Limit:\t\t%d\n", storageBoxType.SubaccountsLimit)
		cmd.Print(util.DescribeDeprecation(storageBoxType))

		// TODO: Skipping prices right now because we have no currency.
		// cmd.Printf("Pricings per Location:\n")
		// for _, price := range storageBoxType.Pricings {
		// 	cmd.Printf("  - Location:\t\t%s\n", price.Location)
		// 	cmd.Printf("    Hourly:\t\t%s\n", util.GrossPrice(price.PriceHourly))
		// 	cmd.Printf("    Monthly:\t\t%s\n", util.GrossPrice(price.PriceMonthly))
		// 	cmd.Printf("    Setup Fee:\t%s\n", util.GrossPrice(price.SetupFee))
		// 	cmd.Printf("\n")
		// }

		return nil
	},
}

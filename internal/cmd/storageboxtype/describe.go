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

		err := loadCurrencyFromAPI(s, storageBoxType)
		if err != nil {
			cmd.PrintErrf("failed to get currency for Storage Box Type prices: %v", err)
		}

		cmd.Printf("Pricings per Location:\n")
		for _, price := range storageBoxType.Pricings {
			cmd.Printf("  - Location:\t%s\n", price.Location)
			cmd.Printf("    Hourly:\t%s\n", util.GrossPrice(price.PriceHourly))
			cmd.Printf("    Monthly:\t%s\n", util.GrossPrice(price.PriceMonthly))
			cmd.Printf("    Setup Fee:\t%s\n", util.GrossPrice(price.SetupFee))
			cmd.Printf("\n")
		}

		return nil
	},
}

func loadCurrencyFromAPI(s state.State, storageBoxType *hcloud.StorageBoxType) error {
	pricing, _, err := s.Client().Pricing().Get(s)
	if err != nil {
		return err
	}

	for i := range storageBoxType.Pricings {
		storageBoxType.Pricings[i].PriceMonthly.Currency = pricing.Currency
		storageBoxType.Pricings[i].PriceMonthly.VATRate = pricing.VATRate

		storageBoxType.Pricings[i].PriceHourly.Currency = pricing.Currency
		storageBoxType.Pricings[i].PriceHourly.VATRate = pricing.VATRate

		storageBoxType.Pricings[i].SetupFee.Currency = pricing.Currency
		storageBoxType.Pricings[i].SetupFee.VATRate = pricing.VATRate
	}

	return nil
}

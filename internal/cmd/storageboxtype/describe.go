package storageboxtype

import (
	"fmt"
	"io"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/experimental"
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
	PrintText: func(s state.State, _ *cobra.Command, out io.Writer, storageBoxType *hcloud.StorageBoxType) error {
		description, err := DescribeStorageBoxType(s, storageBoxType, false)
		if err != nil {
			return err
		}
		fmt.Fprint(out, description)
		return nil
	},
	Experimental: experimental.StorageBoxes,
}

func DescribeStorageBoxType(s state.State, storageBoxType *hcloud.StorageBoxType, short bool) (string, error) {
	var sb strings.Builder

	fmt.Fprintf(&sb, "ID:\t%d\n", storageBoxType.ID)
	fmt.Fprintf(&sb, "Name:\t%s\n", storageBoxType.Name)
	fmt.Fprintf(&sb, "Description:\t%s\n", storageBoxType.Description)
	fmt.Fprintf(&sb, "Size:\t%s\n", humanize.IBytes(uint64(storageBoxType.Size)))
	if storageBoxType.SnapshotLimit != nil {
		fmt.Fprintf(&sb, "Snapshot Limit:\t%d\n", *storageBoxType.SnapshotLimit)
	}
	if storageBoxType.AutomaticSnapshotLimit != nil {
		fmt.Fprintf(&sb, "Automatic Snapshot Limit:\t%d\n", *storageBoxType.AutomaticSnapshotLimit)
	}
	fmt.Fprintf(&sb, "Subaccounts Limit:\t%d\n", storageBoxType.SubaccountsLimit)

	if storageBoxType.IsDeprecated() {
		fmt.Fprintln(&sb)
		fmt.Fprint(&sb, util.DescribeDeprecation(storageBoxType))
	}

	if short {
		return sb.String(), nil
	}

	err := loadCurrencyFromAPI(s, storageBoxType)
	if err != nil {
		return "", fmt.Errorf("failed to get currency for Storage Box Type prices: %w", err)
	}

	fmt.Fprintln(&sb)
	fmt.Fprintf(&sb, "Pricings per Location:\n")
	for i, price := range storageBoxType.Pricings {
		if i > 0 {
			fmt.Fprintln(&sb)
		}
		fmt.Fprintf(&sb, "  - Location:\t%s\n", price.Location)
		fmt.Fprintf(&sb, "    Hourly:\t%s\n", util.GrossPrice(price.PriceHourly))
		fmt.Fprintf(&sb, "    Monthly:\t%s\n", util.GrossPrice(price.PriceMonthly))
		fmt.Fprintf(&sb, "    Setup Fee:\t%s\n", util.GrossPrice(price.SetupFee))
	}

	return sb.String(), nil
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

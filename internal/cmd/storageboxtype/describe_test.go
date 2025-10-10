package storageboxtype_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/storageboxtype"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := storageboxtype.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.StorageBoxTypeClient.EXPECT().
		Get(gomock.Any(), "bx11").
		Return(&hcloud.StorageBoxType{
			ID:                     42,
			Name:                   "bx11",
			Description:            "BX11",
			SnapshotLimit:          hcloud.Ptr(10),
			AutomaticSnapshotLimit: hcloud.Ptr(10),
			SubaccountsLimit:       200,
			Size:                   1073741824,
			Pricings: []hcloud.StorageBoxTypeLocationPricing{
				{
					Location: "fsn1",
					PriceHourly: hcloud.Price{
						Gross: "0.0051",
						Net:   "0.0051",
					},
					PriceMonthly: hcloud.Price{
						Gross: "3.2000",
						Net:   "3.2000",
					},
					SetupFee: hcloud.Price{
						Gross: "0.0000",
						Net:   "0.0000",
					},
				},
				{
					Location: "hel1",
					PriceHourly: hcloud.Price{
						Gross: "0.0051",
						Net:   "0.0051",
					},
					PriceMonthly: hcloud.Price{
						Gross: "3.2000",
						Net:   "3.2000",
					},
					SetupFee: hcloud.Price{
						Gross: "0.0000",
						Net:   "0.0000",
					},
				},
			},
		}, nil, nil)

	fx.Client.PricingClient.EXPECT().
		Get(gomock.Any()).
		Return(hcloud.Pricing{
			Currency: "EUR",
			VATRate:  "1.19",
		}, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"bx11"})

	expOut := `ID:				42
Name:				bx11
Description:			BX11
Size:				1.0 GiB
Snapshot Limit:			10
Automatic Snapshot Limit:	10
Subaccounts Limit:		200
Pricings per Location:
  - Location:	fsn1
    Hourly:	€ 0.0051
    Monthly:	€ 3.2000
    Setup Fee:	€ 0.0000

  - Location:	hel1
    Hourly:	€ 0.0051
    Monthly:	€ 3.2000
    Setup Fee:	€ 0.0000

`

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

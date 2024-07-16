package servertype_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/servertype"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := servertype.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ServerTypeClient.EXPECT().
		Get(gomock.Any(), "cax11").
		Return(&hcloud.ServerType{
			ID:          45,
			Name:        "cax11",
			Description: "CAX11",
			Cores:       2,
			CPUType:     hcloud.CPUTypeShared,
			Memory:      4.0,
			Disk:        40,
			StorageType: hcloud.StorageTypeLocal,
		}, nil, nil)

	fx.Client.PricingClient.EXPECT().
		Get(gomock.Any()).
		Return(hcloud.Pricing{
			ServerTypes: []hcloud.ServerTypePricing{
				// Two server types to test that fullPricingInfo filters for the correct one
				{
					ServerType: &hcloud.ServerType{ID: 1},
					Pricings: []hcloud.ServerTypeLocationPricing{{
						Location: &hcloud.Location{
							Name: "Nuremberg",
						},
						Hourly: hcloud.Price{
							Gross:    "4.0000",
							Currency: "EUR",
						},
						Monthly: hcloud.Price{
							Gross:    "7.0000",
							Currency: "EUR",
						},
					}},
				},
				{
					ServerType: &hcloud.ServerType{ID: 45},
					Pricings: []hcloud.ServerTypeLocationPricing{{
						Location: &hcloud.Location{
							Name: "Falkenstein",
						},
						Hourly: hcloud.Price{
							Gross:    "1.0000",
							Currency: "EUR",
						},
						Monthly: hcloud.Price{
							Gross:    "2.0000",
							Currency: "EUR",
						},
					}},
				},
			},
		}, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"cax11"})

	expOut := `ID:			45
Name:			cax11
Description:		CAX11
Cores:			2
CPU Type:		shared
Architecture:		
Memory:			4.0 GB
Disk:			40 GB
Storage Type:		local
Included Traffic:	0 TB
Pricings per Location:
  - Location:	Falkenstein
    Hourly:	€ 1.0000
    Monthly:	€ 2.0000
`

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

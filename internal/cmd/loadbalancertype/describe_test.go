package loadbalancertype_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/loadbalancertype"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := loadbalancertype.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerTypeClient.EXPECT().
		Get(gomock.Any(), "lb11").
		Return(&hcloud.LoadBalancerType{
			ID:                      123,
			Name:                    "lb11",
			Description:             "LB11",
			MaxServices:             5,
			MaxConnections:          10000,
			MaxTargets:              25,
			MaxAssignedCertificates: 10,
		}, nil, nil)

	fx.Client.PricingClient.EXPECT().
		Get(gomock.Any()).
		Return(hcloud.Pricing{
			LoadBalancerTypes: []hcloud.LoadBalancerTypePricing{
				// Two load balancer types to test that fullPricingInfo filters for the correct one
				{
					LoadBalancerType: &hcloud.LoadBalancerType{ID: 1},
					Pricings: []hcloud.LoadBalancerTypeLocationPricing{{
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
						IncludedTraffic: 6543210,
						PerTBTraffic: hcloud.Price{
							Gross:    "8.0000",
							Currency: "EUR",
						},
					}},
				},
				{
					LoadBalancerType: &hcloud.LoadBalancerType{ID: 123},
					Pricings: []hcloud.LoadBalancerTypeLocationPricing{{
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
						IncludedTraffic: 654321,
						PerTBTraffic: hcloud.Price{
							Gross:    "3.0000",
							Currency: "EUR",
						},
					}},
				},
			},
		}, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{"lb11"})

	expOut := `ID:				123
Name:				lb11
Description:			LB11
Max Services:			5
Max Connections:		10000
Max Targets:			25
Max assigned Certificates:	10
Pricings per Location:
  - Location:		Falkenstein
    Hourly:		€ 1.0000
    Monthly:		€ 2.0000
    Included Traffic:	639 KiB
    Additional Traffic:	€ 3.0000 per TB

`

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

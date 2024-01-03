package loadbalancertype

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := DescribeCmd.CobraCommand(fx.State())
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

	out, _, err := fx.Run(cmd, []string{"lb11"})

	expOut := `ID:				123
Name:				lb11
Description:			LB11
Max Services:			5
Max Connections:		10000
Max Targets:			25
Max assigned Certificates:	10
Pricings per Location:
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

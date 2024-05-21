package loadbalancertype_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/loadbalancertype"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDescribe(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := loadbalancertype.DescribeCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.LoadBalancerType.EXPECT().
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

	out, errOut, err := fx.Run(cmd, []string{"lb11"})

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
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

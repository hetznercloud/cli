package network_test

import (
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/network"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestRemoveSubnet(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := network.RemoveSubnetCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	_, ipRange, _ := net.ParseCIDR("10.0.0.0/24")
	fx.Client.NetworkClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Network{
			ID: 123,
		}, nil, nil)
	fx.Client.NetworkClient.EXPECT().
		DeleteSubnet(gomock.Any(), &hcloud.Network{ID: 123}, hcloud.NetworkDeleteSubnetOpts{
			Subnet: hcloud.NetworkSubnet{
				IPRange: ipRange,
			},
		}).
		Return(&hcloud.Action{ID: 456}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(fx.State(), gomock.Any(), &hcloud.Action{ID: 456}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"123", "--ip-range", "10.0.0.0/24"})

	expOut := "Subnet 10.0.0.0/24 removed from network 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.Empty(t, errOut)
}

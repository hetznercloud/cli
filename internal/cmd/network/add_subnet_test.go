package network_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/network"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestAddSubnet(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := network.AddSubnetCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	_, ipRange, _ := net.ParseCIDR("10.0.0.0/24")
	fx.Client.Network.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Network{
			ID: 123,
		}, nil, nil)
	fx.Client.Network.EXPECT().
		AddSubnet(gomock.Any(), &hcloud.Network{ID: 123}, hcloud.NetworkAddSubnetOpts{
			Subnet: hcloud.NetworkSubnet{
				Type:        "cloud",
				NetworkZone: "eu-central",
				IPRange:     ipRange,
			},
		}).
		Return(&hcloud.Action{ID: 456}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 456}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"123", "--type", "cloud", "--network-zone", "eu-central", "--ip-range", "10.0.0.0/24"})

	expOut := "Subnet added to network 123\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.Empty(t, errOut)
}

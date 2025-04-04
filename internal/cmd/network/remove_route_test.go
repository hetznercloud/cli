package network_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/network"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestRemoveRoute(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := network.RemoveRouteCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	_, ipRange, _ := net.ParseCIDR("10.0.0.0/24")
	fx.Client.NetworkClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Network{
			ID: 123,
		}, nil, nil)
	fx.Client.NetworkClient.EXPECT().
		DeleteRoute(gomock.Any(), &hcloud.Network{ID: 123}, hcloud.NetworkDeleteRouteOpts{
			Route: hcloud.NetworkRoute{
				Destination: ipRange,
				Gateway:     net.IPv4(10, 0, 0, 1),
			},
		}).
		Return(&hcloud.Action{ID: 456}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(fx.State(), gomock.Any(), &hcloud.Action{ID: 456}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"123", "--destination", "10.0.0.0/24", "--gateway", "10.0.0.1"})

	expOut := "Route removed from Network 123\n"

	require.NoError(t, err)
	assert.Equal(t, expOut, out)
	assert.Empty(t, errOut)
}

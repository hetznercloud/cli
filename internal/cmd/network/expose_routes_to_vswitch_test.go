package network_test

import (
	_ "embed"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/network"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestExposeRoutesToVSwitchEnable(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := network.ExposeRoutesToVSwitchCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	n := &hcloud.Network{ID: 123, Name: "myNetwork"}

	fx.Client.NetworkClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(n, nil, nil)
	fx.Client.NetworkClient.EXPECT().
		Update(gomock.Any(), n, hcloud.NetworkUpdateOpts{
			ExposeRoutesToVSwitch: hcloud.Ptr(true),
		}).
		Return(n, nil, nil)

	args := []string{"123"}
	out, errOut, err := fx.Run(cmd, args)

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Exposing routes to connected vSwitch of network myNetwork enabled\n", out)
}

func TestExposeRoutesToVSwitchDisable(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := network.ExposeRoutesToVSwitchCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	n := &hcloud.Network{ID: 123, Name: "myNetwork"}

	fx.Client.NetworkClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(n, nil, nil)
	fx.Client.NetworkClient.EXPECT().
		Update(gomock.Any(), n, hcloud.NetworkUpdateOpts{
			ExposeRoutesToVSwitch: hcloud.Ptr(false),
		}).
		Return(n, nil, nil)

	args := []string{"123", "--disable"}
	out, errOut, err := fx.Run(cmd, args)

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Exposing routes to connected vSwitch of network myNetwork disabled\n", out)
}

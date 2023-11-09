package network

import (
	"context"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestCreate(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	_, ipRange, _ := net.ParseCIDR("10.0.0.0/24")
	fx.Client.NetworkClient.EXPECT().
		Create(gomock.Any(), hcloud.NetworkCreateOpts{
			Name:    "myNetwork",
			IPRange: ipRange,
			Labels:  make(map[string]string),
		}).
		Return(&hcloud.Network{
			ID:      123,
			Name:    "myNetwork",
			IPRange: ipRange,
		}, nil, nil)

	out, _, err := fx.Run(cmd, []string{"--name", "myNetwork", "--ip-range", "10.0.0.0/24"})

	expOut := "Network 123 created\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestCreateProtection(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := CreateCmd.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	_, ipRange, _ := net.ParseCIDR("10.0.0.0/24")
	network := &hcloud.Network{
		ID:      123,
		Name:    "myNetwork",
		IPRange: ipRange,
	}

	fx.Client.NetworkClient.EXPECT().
		Create(gomock.Any(), hcloud.NetworkCreateOpts{
			Name:    "myNetwork",
			IPRange: ipRange,
			Labels:  make(map[string]string),
		}).
		Return(network, nil, nil)
	fx.Client.NetworkClient.EXPECT().
		ChangeProtection(gomock.Any(), network, hcloud.NetworkChangeProtectionOpts{
			Delete: hcloud.Ptr(true),
		}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), &hcloud.Action{ID: 123}).Return(nil)

	out, _, err := fx.Run(cmd, []string{"--name", "myNetwork", "--ip-range", "10.0.0.0/24", "--enable-protection", "delete"})

	expOut := `Network 123 created
Resource protection enabled for network 123
`

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

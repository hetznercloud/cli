package network_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/network"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDisableProtection(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := network.DisableProtectionCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	n := &hcloud.Network{ID: 123, Name: "myNetwork"}

	fx.Client.NetworkClient.EXPECT().
		Get(gomock.Any(), "myNetwork").
		Return(n, nil, nil)
	fx.Client.NetworkClient.EXPECT().
		ChangeProtection(gomock.Any(), n, hcloud.NetworkChangeProtectionOpts{
			Delete: hcloud.Ptr(false),
		}).
		Return(&hcloud.Action{ID: 789}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)

	args := []string{"myNetwork", "delete"}
	out, errOut, err := fx.Run(cmd, args)

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Resource protection disabled for network 123\n", out)
}

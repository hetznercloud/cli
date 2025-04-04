package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cli"
	"github.com/hetznercloud/cli/internal/state/config"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestOptionFlagParsing(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := cli.NewRootCommand(fx.State())
	fx.ExpectEnsureToken()

	net := &hcloud.Network{ID: 1, Name: "foo"}
	srv := &hcloud.Server{ID: 2, Name: "bar"}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "bar").
		Return(srv, nil, nil)
	fx.Client.NetworkClient.EXPECT().
		Get(gomock.Any(), "foo").
		Return(net, nil, nil)
	fx.Client.ServerClient.EXPECT().
		AttachToNetwork(gomock.Any(), srv, hcloud.ServerAttachToNetworkOpts{Network: net}).
		Return(&hcloud.Action{ID: 3}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 3}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"server", "attach-to-network", "--network", "foo", "bar", "--debug"})

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Server 2 attached to Network 1\n", out)

	val, err := config.OptionDebug.Get(fx.State().Config())
	require.NoError(t, err)
	assert.True(t, val)
}

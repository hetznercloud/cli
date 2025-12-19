package server_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestRequestConsole(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.RequestConsoleCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := &hcloud.Server{ID: 123, Name: "my-server"}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "my-server").
		Return(srv, nil, nil)
	fx.Client.ServerClient.EXPECT().
		RequestConsole(gomock.Any(), srv).
		Return(hcloud.ServerRequestConsoleResult{
			Action:   &hcloud.Action{ID: 789},
			Password: "root-password",
			WSSURL:   "wss://console.hetzner.com/?token=123",
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)

	args := []string{"my-server"}
	out, errOut, err := fx.Run(cmd, args)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Console for Server 123:\nWebSocket URL: wss://console.hetzner.com/?token=123\nVNC Password: root-password\n", out)
}

func TestRequestConsoleJSON(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.RequestConsoleCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := &hcloud.Server{ID: 123, Name: "my-server"}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "my-server").
		Return(srv, nil, nil)
	fx.Client.ServerClient.EXPECT().
		RequestConsole(gomock.Any(), srv).
		Return(hcloud.ServerRequestConsoleResult{
			Action:   &hcloud.Action{ID: 789},
			Password: "root-password",
			WSSURL:   "wss://console.hetzner.com/?token=123",
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)

	args := []string{"my-server", "-o=json"}
	out, errOut, err := fx.Run(cmd, args)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.JSONEq(t, `{"wss_url": "wss://console.hetzner.com/?token=123", "password": "root-password"}`, out)
}

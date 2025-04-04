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

func TestShutdown(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.ShutdownCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := hcloud.Server{
		ID:     42,
		Name:   "my server",
		Status: hcloud.ServerStatusRunning,
	}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), srv.Name).
		Return(&srv, nil, nil)

	fx.Client.ServerClient.EXPECT().
		Shutdown(gomock.Any(), &srv)
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), gomock.Any(), nil)

	out, errOut, err := fx.Run(cmd, []string{srv.Name})

	expOut := "Sent shutdown signal to Server 42\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestShutdownWait(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.ShutdownCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := hcloud.Server{
		ID:     42,
		Name:   "my server",
		Status: hcloud.ServerStatusRunning,
	}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), srv.Name).
		Return(&srv, nil, nil)

	fx.Client.ServerClient.EXPECT().
		Shutdown(gomock.Any(), &srv)
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), gomock.Any(), nil)

	fx.Client.ServerClient.EXPECT().
		GetByID(gomock.Any(), srv.ID).
		Return(&srv, nil, nil).
		Return(&srv, nil, nil).
		Return(&hcloud.Server{ID: srv.ID, Name: srv.Name, Status: hcloud.ServerStatusOff}, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{srv.Name, "--wait"})

	expOut := "Sent shutdown signal to Server 42\nServer 42 shut down\n"
	expErrOut := "Waiting for Server to shut down (server: 42) ...\nWaiting for Server to shut down (server: 42) ... done\n"

	require.NoError(t, err)
	assert.Equal(t, expErrOut, errOut)
	assert.Equal(t, expOut, out)
}

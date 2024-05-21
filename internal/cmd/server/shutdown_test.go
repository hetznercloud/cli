package server_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	fx.Client.Server.EXPECT().
		Get(gomock.Any(), srv.Name).
		Return(&srv, nil, nil)

	fx.Client.Server.EXPECT().
		Shutdown(gomock.Any(), &srv)
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), gomock.Any(), nil)

	out, errOut, err := fx.Run(cmd, []string{srv.Name})

	expOut := "Sent shutdown signal to server 42\n"

	assert.NoError(t, err)
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

	fx.Client.Server.EXPECT().
		Get(gomock.Any(), srv.Name).
		Return(&srv, nil, nil)

	fx.Client.Server.EXPECT().
		Shutdown(gomock.Any(), &srv)
	fx.ActionWaiter.EXPECT().WaitForActions(gomock.Any(), gomock.Any(), nil)

	fx.Client.Server.EXPECT().
		GetByID(gomock.Any(), srv.ID).
		Return(&srv, nil, nil).
		Return(&srv, nil, nil).
		Return(&hcloud.Server{ID: srv.ID, Name: srv.Name, Status: hcloud.ServerStatusOff}, nil, nil)

	out, errOut, err := fx.Run(cmd, []string{srv.Name, "--wait"})

	expOut := "Sent shutdown signal to server 42\nServer 42 shut down\n"
	expErrOut := "Waiting for server to shut down (server: 42) ...\nWaiting for server to shut down (server: 42) ... done\n"

	assert.NoError(t, err)
	assert.Equal(t, expErrOut, errOut)
	assert.Equal(t, expOut, out)
}

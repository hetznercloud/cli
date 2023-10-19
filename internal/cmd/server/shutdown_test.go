package server

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShutdown(t *testing.T) {

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := ShutdownCommand.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	var (
		server = hcloud.Server{
			ID:     42,
			Name:   "my server",
			Status: hcloud.ServerStatusRunning,
		}
	)

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), server.Name).
		Return(&server, nil, nil)

	fx.Client.ServerClient.EXPECT().
		Shutdown(gomock.Any(), &server)
	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), nil)

	out, err := fx.Run(cmd, []string{server.Name})

	expOut := "Sent shutdown signal to server 42\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestShutdownWait(t *testing.T) {

	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := ShutdownCommand.CobraCommand(
		context.Background(),
		fx.Client,
		fx.TokenEnsurer,
		fx.ActionWaiter)
	fx.ExpectEnsureToken()

	var (
		server = hcloud.Server{
			ID:     42,
			Name:   "my server",
			Status: hcloud.ServerStatusRunning,
		}
	)

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), server.Name).
		Return(&server, nil, nil)

	fx.Client.ServerClient.EXPECT().
		Shutdown(gomock.Any(), &server)
	fx.ActionWaiter.EXPECT().ActionProgress(gomock.Any(), nil)

	fx.Client.ServerClient.EXPECT().
		GetByID(gomock.Any(), server.ID).
		Return(&server, nil, nil).
		Return(&server, nil, nil).
		Return(&hcloud.Server{ID: server.ID, Name: server.Name, Status: hcloud.ServerStatusOff}, nil, nil)

	out, err := fx.Run(cmd, []string{server.Name, "--wait"})

	expOut := "Sent shutdown signal to server 42\nWaiting for server to shut down ... done\nServer 42 shut down\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

package server_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestRebuild(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.RebuildCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := &hcloud.Server{ID: 123, Name: "my-server", ServerType: &hcloud.ServerType{Architecture: hcloud.ArchitectureARM}}
	img := &hcloud.Image{ID: 456, Name: "ubuntu-22.04"}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "my-server").
		Return(srv, nil, nil)
	fx.Client.ImageClient.EXPECT().
		GetForArchitecture(gomock.Any(), "ubuntu-22.04", hcloud.ArchitectureARM).
		Return(img, nil, nil)
	fx.Client.ServerClient.EXPECT().
		RebuildWithResult(gomock.Any(), srv, hcloud.ServerRebuildOpts{Image: img}).
		Return(hcloud.ServerRebuildResult{
			Action:       &hcloud.Action{ID: 789},
			RootPassword: "root-password",
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)

	args := []string{"my-server", "--image", "ubuntu-22.04"}
	out, errOut, err := fx.Run(cmd, args)

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Server 123 rebuilt with image ubuntu-22.04\nRoot password: root-password\n", out)
}

func TestRebuildDeprecated(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.RebuildCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	cmd.SilenceUsage = true // Silence usage output for this test; usually handled by root command

	srv := &hcloud.Server{ID: 123, Name: "my-server", ServerType: &hcloud.ServerType{Architecture: hcloud.ArchitectureARM}}
	img := &hcloud.Image{ID: 456, Name: "ubuntu-22.04", Deprecated: time.Date(2036, 5, 20, 0, 0, 0, 0, time.UTC)}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "my-server").
		Return(srv, nil, nil)
	fx.Client.ImageClient.EXPECT().
		GetForArchitecture(gomock.Any(), "ubuntu-22.04", hcloud.ArchitectureARM).
		Return(img, nil, nil)

	args := []string{"my-server", "--image", "ubuntu-22.04"}
	out, errOut, err := fx.Run(cmd, args)

	errorMsg := "image ubuntu-22.04 is deprecated, please use --allow-deprecated-image to create a server with this image. It will continue to be available until 2036-08-20"

	require.Error(t, err, errorMsg)
	assert.Equal(t, "Error: "+errorMsg+"\n", errOut)
	assert.Empty(t, out)
}

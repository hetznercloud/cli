package server_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/server"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestCreateImage(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := server.CreateImageCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	srv := &hcloud.Server{ID: 123, Name: "my-server"}

	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "my-server").
		Return(srv, nil, nil)
	fx.Client.ServerClient.EXPECT().
		CreateImage(gomock.Any(), srv, &hcloud.ServerCreateImageOpts{
			Type:        hcloud.ImageTypeSnapshot,
			Description: hcloud.Ptr("my-snapshot"),
			Labels:      make(map[string]string),
		}).
		Return(hcloud.ServerCreateImageResult{
			Action: &hcloud.Action{ID: 789},
			Image:  &hcloud.Image{ID: 456},
		}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)

	args := []string{"my-server", "--type", "snapshot", "--description", "my-snapshot"}
	out, errOut, err := fx.Run(cmd, args)

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Image 456 created from server 123\n", out)
}

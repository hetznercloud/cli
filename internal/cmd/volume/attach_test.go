package volume_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/volume"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestAttach(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := volume.AttachCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	v := &hcloud.Volume{ID: 123}
	srv := &hcloud.Server{ID: 456, Name: "myServer"}

	fx.Client.VolumeClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(v, nil, nil)
	fx.Client.ServerClient.EXPECT().
		Get(gomock.Any(), "456").
		Return(srv, nil, nil)
	fx.Client.VolumeClient.EXPECT().
		AttachWithOpts(gomock.Any(), v, hcloud.VolumeAttachOpts{
			Server:    srv,
			Automount: hcloud.Ptr(false),
		}).Return(&hcloud.Action{ID: 321}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 321}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"123", "--server", "456"})

	expOut := "Volume 123 attached to server myServer\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

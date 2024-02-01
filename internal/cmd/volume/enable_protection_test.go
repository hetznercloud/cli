package volume_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/volume"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestEnableProtection(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := volume.EnableProtectionCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	v := &hcloud.Volume{ID: 123, Name: "myVolume"}

	fx.Client.VolumeClient.EXPECT().
		Get(gomock.Any(), "myVolume").
		Return(v, nil, nil)
	fx.Client.VolumeClient.EXPECT().
		ChangeProtection(gomock.Any(), v, hcloud.VolumeChangeProtectionOpts{
			Delete: hcloud.Ptr(true),
		}).
		Return(&hcloud.Action{ID: 789}, nil, nil)
	fx.ActionWaiter.EXPECT().
		ActionProgress(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 789}).
		Return(nil)

	args := []string{"myVolume", "delete"}
	out, errOut, err := fx.Run(cmd, args)

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, "Resource protection enabled for volume 123\n", out)
}

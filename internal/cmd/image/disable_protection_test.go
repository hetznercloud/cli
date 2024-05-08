package image_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/cmd/image"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDisableProtection(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := image.DisableProtectionCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ImageClient.EXPECT().
		ChangeProtection(gomock.Any(), &hcloud.Image{ID: 123}, hcloud.ImageChangeProtectionOpts{
			Delete: hcloud.Ptr(false),
		}).
		Return(&hcloud.Action{ID: 123}, nil, nil)
	fx.ActionWaiter.EXPECT().
		WaitForActions(gomock.Any(), gomock.Any(), &hcloud.Action{ID: 123}).
		Return(nil)

	out, errOut, err := fx.Run(cmd, []string{"123", "delete"})

	expOut := "Resource protection disabled for image 123\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

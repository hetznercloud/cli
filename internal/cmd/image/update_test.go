package image_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/hetznercloud/cli/internal/cmd/image"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestUpdateDescription(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := image.UpdateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.Image.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Image{ID: 123}, nil, nil)
	fx.Client.Image.EXPECT().
		Update(gomock.Any(), &hcloud.Image{ID: 123}, hcloud.ImageUpdateOpts{
			Description: hcloud.Ptr("new-description"),
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "--description", "new-description"})

	expOut := "Image 123 updated\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestUpdateType(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := image.UpdateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.Image.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Image{ID: 123}, nil, nil)
	fx.Client.Image.EXPECT().
		Update(gomock.Any(), &hcloud.Image{ID: 123}, hcloud.ImageUpdateOpts{
			Description: hcloud.Ptr(""),
			Type:        hcloud.ImageTypeSnapshot,
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "--type", "snapshot"})

	expOut := "Image 123 updated\n"

	assert.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

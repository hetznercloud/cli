package image_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hetznercloud/cli/internal/cmd/image"
	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestUpdateDescription(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := image.UpdateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ImageClient.EXPECT().
		GetByID(gomock.Any(), int64(123)).
		Return(&hcloud.Image{ID: 123}, nil, nil)
	fx.Client.ImageClient.EXPECT().
		Update(gomock.Any(), &hcloud.Image{ID: 123}, hcloud.ImageUpdateOpts{
			Description: hcloud.Ptr("new-description"),
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "--description", "new-description"})

	expOut := "Image 123 updated\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

func TestUpdateType(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := image.UpdateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ImageClient.EXPECT().
		GetByID(gomock.Any(), int64(123)).
		Return(&hcloud.Image{ID: 123}, nil, nil)
	fx.Client.ImageClient.EXPECT().
		Update(gomock.Any(), &hcloud.Image{ID: 123}, hcloud.ImageUpdateOpts{
			Description: hcloud.Ptr(""),
			Type:        hcloud.ImageTypeSnapshot,
		})

	out, errOut, err := fx.Run(cmd, []string{"123", "--type", "snapshot"})

	expOut := "Image 123 updated\n"

	require.NoError(t, err)
	assert.Empty(t, errOut)
	assert.Equal(t, expOut, out)
}

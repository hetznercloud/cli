package image

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/cli/internal/testutil"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestUpdateDescription(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := UpdateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ImageClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Image{ID: 123}, nil, nil)
	fx.Client.ImageClient.EXPECT().
		Update(gomock.Any(), &hcloud.Image{ID: 123}, hcloud.ImageUpdateOpts{
			Description: hcloud.Ptr("new-description"),
		})

	out, _, err := fx.Run(cmd, []string{"123", "--description", "new-description"})

	expOut := "Image 123 updated\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}

func TestUpdateType(t *testing.T) {
	fx := testutil.NewFixture(t)
	defer fx.Finish()

	cmd := UpdateCmd.CobraCommand(fx.State())
	fx.ExpectEnsureToken()

	fx.Client.ImageClient.EXPECT().
		Get(gomock.Any(), "123").
		Return(&hcloud.Image{ID: 123}, nil, nil)
	fx.Client.ImageClient.EXPECT().
		Update(gomock.Any(), &hcloud.Image{ID: 123}, hcloud.ImageUpdateOpts{
			Description: hcloud.Ptr(""),
			Type:        hcloud.ImageTypeSnapshot,
		})

	out, _, err := fx.Run(cmd, []string{"123", "--type", "snapshot"})

	expOut := "Image 123 updated\n"

	assert.NoError(t, err)
	assert.Equal(t, expOut, out)
}
